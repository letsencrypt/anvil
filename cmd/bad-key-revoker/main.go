package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	netmail "net/mail"
	"os"
	"strings"
	"time"

	"github.com/letsencrypt/boulder/cmd"
	"github.com/letsencrypt/boulder/core"
	corepb "github.com/letsencrypt/boulder/core/proto"
	"github.com/letsencrypt/boulder/db"
	bgrpc "github.com/letsencrypt/boulder/grpc"
	"github.com/letsencrypt/boulder/log"
	"github.com/letsencrypt/boulder/mail"
	rapb "github.com/letsencrypt/boulder/ra/proto"
	"github.com/letsencrypt/boulder/sa"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/crypto/ocsp"

	"google.golang.org/grpc"
)

var keysProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "bad_keys_processed",
	Help: "A counter of blockedKeys rows processed labelled by processing state",
}, []string{"state"})
var certsRevoked = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bad_keys_certs_revoked",
	Help: "A counter of certificates associated with rows in blockedKeys that have been revoked",
})
var mailErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "bad_keys_mail_errors",
	Help: "A counter of email send errors",
})

// revoker is an interface used to reduce the scope of a RA gRPC client
// to only the single method we need to use, this makes testing significantly
// simpler
type revoker interface {
	AdministrativelyRevokeCertificate(ctx context.Context, in *rapb.AdministrativelyRevokeCertificateRequest, opts ...grpc.CallOption) (*corepb.Empty, error)
}

type badKeyRevoker struct {
	dbMap           *db.WrappedMap
	maxRevocations  int
	serialBatchSize int
	raClient        revoker
	mailer          mail.Mailer
	emailSubject    string
	emailTemplate   *template.Template
	logger          log.Logger
}

// uncheckedBlockedKey represents a row in the blockedKeys table
type uncheckedBlockedKey struct {
	KeyHash   []byte
	RevokedBy int64
}

func (ubk uncheckedBlockedKey) String() string {
	return fmt.Sprintf("[revokedBy: %d, keyHash: %x]",
		ubk.RevokedBy, ubk.KeyHash)
}

func (bkr *badKeyRevoker) selectUncheckedKey() (uncheckedBlockedKey, error) {
	var row uncheckedBlockedKey
	err := bkr.dbMap.SelectOne(
		&row,
		`SELECT keyHash, revokedBy
		FROM blockedKeys
		WHERE extantCertificatesChecked = false
		LIMIT 1`,
	)
	return row, err
}

// unrevokedCertificate represents a yet to be revoked certificate
type unrevokedCertificate struct {
	ID             int
	Serial         string
	DER            []byte
	RegistrationID int64
}

// findUnrevoked looks for all unexpired, currently valid certificates which have a specific SPKI hash,
// by looking first at the keyHashToSerial table and then the certificateStatus and certificates tables.
// If the number of certificates it finds is larger than bkr.maxRevocations it'll error out.
func (bkr *badKeyRevoker) findUnrevoked(unchecked uncheckedBlockedKey) ([]unrevokedCertificate, error) {
	var unrevokedCerts []unrevokedCertificate
	initialID := 0
	for {
		var batch []struct {
			ID         int
			CertSerial string
		}
		_, err := bkr.dbMap.Select(
			&batch,
			"SELECT id, certserial FROM keyHashToSerial WHERE keyHash = ? AND id > ? ORDER BY id LIMIT ?",
			unchecked.KeyHash,
			initialID,
			bkr.serialBatchSize,
		)
		if err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		initialID = batch[len(batch)-1].ID
		for _, serial := range batch {
			var unrevokedCert unrevokedCertificate
			err = bkr.dbMap.SelectOne(
				&unrevokedCert,
				`SELECT cs.id, cs.serial, c.registrationID, c.der
				FROM certificateStatus AS cs
				JOIN certificates AS c
				ON cs.serial = c.serial
				WHERE cs.serial = ? AND cs.isExpired = false AND cs.status != ?`,
				serial.CertSerial,
				string(core.StatusRevoked),
			)
			if err != nil {
				if db.IsNoRows(err) {
					continue
				}
				return nil, err
			}
			unrevokedCerts = append(unrevokedCerts, unrevokedCert)
		}
	}
	if len(unrevokedCerts) > bkr.maxRevocations {
		return nil, fmt.Errorf("too many certificates to revoke associated with %x: got %d, max %d", unchecked.KeyHash, len(unrevokedCerts), bkr.maxRevocations)
	}
	return unrevokedCerts, nil
}

// markRowChecked updates a row in the blockedKeys table to mark a keyHash
// as having been checked for extant unrevoked certificates.
func (bkr *badKeyRevoker) markRowChecked(unchecked uncheckedBlockedKey) error {
	_, err := bkr.dbMap.Exec("UPDATE blockedKeys SET extantCertificatesChecked = true WHERE keyHash = ?", unchecked.KeyHash)
	return err
}

// resolveContacts builds a map of id -> email addresses
func (bkr *badKeyRevoker) resolveContacts(ids []int64) (map[int64][]string, error) {
	idToEmail := map[int64][]string{}
	for _, id := range ids {
		var emails struct {
			Contact []string
		}
		err := bkr.dbMap.SelectOne(&emails, "SELECT contact FROM registrations WHERE id = ?", id)
		if err != nil {
			if db.IsNoRows(err) {
				continue
			}
			return nil, err
		}
		if len(emails.Contact) != 0 {
			for _, email := range emails.Contact {
				idToEmail[id] = append(idToEmail[id], strings.TrimPrefix(email, "mailto:"))
			}
		}
	}
	return idToEmail, nil
}

var maxSerials = 100

// sendMessage sends a single email to the provided address with the revoked
// serials
func (bkr *badKeyRevoker) sendMessage(addr string, serials []string) error {
	err := bkr.mailer.Connect()
	if err != nil {
		return err
	}
	defer func() {
		_ = bkr.mailer.Close()
	}()
	mutSerials := make([]string, len(serials))
	copy(mutSerials, serials)
	if len(mutSerials) > maxSerials {
		more := len(mutSerials) - maxSerials
		mutSerials = mutSerials[:maxSerials]
		mutSerials = append(mutSerials, fmt.Sprintf("and %d more certificates.", more))
	}
	message := bytes.NewBuffer(nil)
	err = bkr.emailTemplate.Execute(message, mutSerials)
	if err != nil {
		return err
	}
	err = bkr.mailer.SendMail([]string{addr}, bkr.emailSubject, message.String())
	if err != nil {
		return err
	}
	return nil
}

var keyCompromiseCode = int64(ocsp.KeyCompromise)
var revokerName = "bad-key-revoker"

// revokeCerts revokes all the certificates associated with a particular key hash and sends
// emails to the users that issued the certificates. Emails are not sent to the user which
// requested revocation of the original certificate which marked the key as compromised.
func (bkr *badKeyRevoker) revokeCerts(revokerEmails []string, emailToCerts map[string][]unrevokedCertificate) error {
	revokerEmailsMap := map[string]bool{}
	for _, email := range revokerEmails {
		revokerEmailsMap[email] = true
	}

	alreadyRevoked := map[int]bool{}
	for email, certs := range emailToCerts {
		var revokedSerials []string
		for _, cert := range certs {
			revokedSerials = append(revokedSerials, cert.Serial)
			if alreadyRevoked[cert.ID] {
				continue
			}
			_, err := bkr.raClient.AdministrativelyRevokeCertificate(context.Background(), &rapb.AdministrativelyRevokeCertificateRequest{
				Cert:      cert.DER,
				Code:      &keyCompromiseCode,
				AdminName: &revokerName,
			})
			if err != nil {
				return err
			}
			certsRevoked.Inc()
			alreadyRevoked[cert.ID] = true
		}
		// don't send emails to the person who revoked the certificate
		if revokerEmailsMap[email] || email == "" {
			continue
		}
		err := bkr.sendMessage(email, revokedSerials)
		if err != nil {
			mailErrors.Inc()
			bkr.logger.Errf("failed to send message to %q: %s", email, err)
			continue
		}
	}
	return nil
}

// invoke processes a single key in the blockedKeys table and returns whether
// there were any rows to process or not.
func (bkr *badKeyRevoker) invoke() (bool, error) {
	// select a row to process
	unchecked, err := bkr.selectUncheckedKey()
	if err != nil {
		if db.IsNoRows(err) {
			return true, nil
		}
		return false, err
	}
	bkr.logger.AuditInfo(fmt.Sprintf("found unchecked block key to work on: %s", unchecked))

	// select all unrevoked, unexpired serials associated with the blocked key hash
	unrevokedCerts, err := bkr.findUnrevoked(unchecked)
	if err != nil {
		bkr.logger.AuditInfo(fmt.Sprintf("finding unrevoked certificates related to %s: %s",
			unchecked, err))
		return false, err
	}
	if len(unrevokedCerts) == 0 {
		bkr.logger.AuditInfo(fmt.Sprintf("found no certificates that need revoking related to %s, marking row as checked", unchecked))
		// mark row as checked
		err = bkr.markRowChecked(unchecked)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// build a map of registration ID -> certificates, and collect a
	// list of unique registration IDs
	ownedBy := map[int64][]unrevokedCertificate{}
	var ids []int64
	for _, cert := range unrevokedCerts {
		if ownedBy[cert.RegistrationID] == nil {
			ids = append(ids, cert.RegistrationID)
		}
		ownedBy[cert.RegistrationID] = append(ownedBy[cert.RegistrationID], cert)
	}
	// if the account that revoked the original certificate isn't an owner of any
	// extant certificates, still add them to ids so that we can resolve their
	// email and avoid sending emails later. If RevokedBy == 0 it was a row
	// inserted by admin-revoker with a dummy ID, since there won't be a registration
	// to look up, don't bother adding it to ids.
	if _, present := ownedBy[unchecked.RevokedBy]; !present && unchecked.RevokedBy != 0 {
		ids = append(ids, unchecked.RevokedBy)
	}
	// get contact addresses for the list of IDs
	idToEmails, err := bkr.resolveContacts(ids)
	if err != nil {
		return false, err
	}

	// build a map of email -> certificates, this de-duplicates accounts with
	// the same email addresses
	emailsToCerts := map[string][]unrevokedCertificate{}
	for id, emails := range idToEmails {
		for _, email := range emails {
			emailsToCerts[email] = append(emailsToCerts[email], ownedBy[id]...)
		}
	}

	revokerEmails := idToEmails[unchecked.RevokedBy]
	bkr.logger.AuditInfo(fmt.Sprintf("revoking certs. revoked emails=%v, emailsToCerts=%v",
		revokerEmails, emailsToCerts))

	// revoke each certificate and send emails to their owners
	err = bkr.revokeCerts(idToEmails[unchecked.RevokedBy], emailsToCerts)
	if err != nil {
		return false, err
	}

	// mark the key as checked
	err = bkr.markRowChecked(unchecked)
	if err != nil {
		return false, err
	}
	return false, nil
}

func main() {
	var config struct {
		BadKeyRevoker struct {
			cmd.DBConfig
			DebugAddr string

			TLS       cmd.TLSConfig
			RAService *cmd.GRPCClientConfig

			// MaximumRevocations specifies the maximum number of certificates associated with
			// a key hash that bad-key-revoker will attempt to revoke. If the number of certificates
			// is higher than MaximumRevocations bad-key-revoker will error out and refuse to
			// progress until this is addressed.
			MaximumRevocations int
			// FindCertificatesBatchSize specifies the maximum number of serials to select from the
			// keyHashToSerial table at once
			FindCertificatesBatchSize int

			// Interval specifies how long bad-key-revoker should sleep between attempting to find
			// blockedKeys rows to process when there is no work to do
			Interval cmd.ConfigDuration

			Mailer struct {
				cmd.SMTPConfig
				// Path to a file containing a list of trusted root certificates for use
				// during the SMTP connection (as opposed to the gRPC connections).
				SMTPTrustedRootFile string

				From          string
				EmailSubject  string
				EmailTemplate string
			}
		}

		Syslog cmd.SyslogConfig
	}
	configPath := flag.String("config", "", "File path to the configuration file for this service")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	err := cmd.ReadConfigFile(*configPath, &config)
	cmd.FailOnError(err, "Failed reading config file")

	scope, logger := cmd.StatsAndLogging(config.Syslog, config.BadKeyRevoker.DebugAddr)
	clk := cmd.Clock()

	scope.MustRegister(keysProcessed)
	scope.MustRegister(certsRevoked)
	scope.MustRegister(mailErrors)

	dbURL, err := config.BadKeyRevoker.DBConfig.URL()
	cmd.FailOnError(err, "Couldn't load DB URL")
	dbMap, err := sa.NewDbMap(dbURL, config.BadKeyRevoker.DBConfig.MaxDBConns)
	cmd.FailOnError(err, "Could not connect to database")
	sa.SetSQLDebug(dbMap, logger)
	sa.InitDBMetrics(dbMap, scope)

	tlsConfig, err := config.BadKeyRevoker.TLS.Load()
	cmd.FailOnError(err, "TLS config")

	clientMetrics := bgrpc.NewClientMetrics(scope)
	conn, err := bgrpc.ClientSetup(config.BadKeyRevoker.RAService, tlsConfig, clientMetrics, clk)
	cmd.FailOnError(err, "Failed to load credentials and create gRPC connection to RA")
	rac := rapb.NewRegistrationAuthorityClient(conn)

	var smtpRoots *x509.CertPool
	if config.BadKeyRevoker.Mailer.SMTPTrustedRootFile != "" {
		pem, err := ioutil.ReadFile(config.BadKeyRevoker.Mailer.SMTPTrustedRootFile)
		cmd.FailOnError(err, "Loading trusted roots file")
		smtpRoots = x509.NewCertPool()
		if !smtpRoots.AppendCertsFromPEM(pem) {
			cmd.FailOnError(nil, "Failed to parse root certs PEM")
		}
	}

	fromAddress, err := netmail.ParseAddress(config.BadKeyRevoker.Mailer.From)
	cmd.FailOnError(err, fmt.Sprintf("Could not parse from address: %s", config.BadKeyRevoker.Mailer.From))

	smtpPassword, err := config.BadKeyRevoker.Mailer.PasswordConfig.Pass()
	cmd.FailOnError(err, "Failed to load SMTP password")
	mailClient := mail.New(
		config.BadKeyRevoker.Mailer.Server,
		config.BadKeyRevoker.Mailer.Port,
		config.BadKeyRevoker.Mailer.Username,
		smtpPassword,
		smtpRoots,
		*fromAddress,
		logger,
		scope,
		1*time.Second,    // reconnection base backoff
		5*60*time.Second, // reconnection maximum backoff
	)

	if config.BadKeyRevoker.Mailer.EmailSubject == "" {
		cmd.Fail("BadKeyRevoker.Mailer.EmailSubject must be populated")
	}
	templateBytes, err := ioutil.ReadFile(config.BadKeyRevoker.Mailer.EmailTemplate)
	cmd.FailOnError(err, fmt.Sprintf("failed to read email template %q: %s", config.BadKeyRevoker.Mailer.EmailTemplate, err))
	emailTemplate, err := template.New("email").Parse(string(templateBytes))
	cmd.FailOnError(err, fmt.Sprintf("failed to parse email template %q: %s", config.BadKeyRevoker.Mailer.EmailTemplate, err))

	bkr := &badKeyRevoker{
		dbMap:           dbMap,
		maxRevocations:  config.BadKeyRevoker.MaximumRevocations,
		serialBatchSize: config.BadKeyRevoker.FindCertificatesBatchSize,
		raClient:        rac,
		mailer:          mailClient,
		emailSubject:    config.BadKeyRevoker.Mailer.EmailSubject,
		emailTemplate:   emailTemplate,
		logger:          logger,
	}
	for {
		noWork, err := bkr.invoke()
		if err != nil {
			keysProcessed.WithLabelValues("error").Inc()
			logger.AuditErrf("failed to process blockedKeys row: %s", err)
			continue
		}
		if noWork {
			logger.AuditInfo(fmt.Sprintf(
				"No work to do. Sleeping for %s", config.BadKeyRevoker.Interval.Duration))
			time.Sleep(config.BadKeyRevoker.Interval.Duration)
		} else {
			keysProcessed.WithLabelValues("success").Inc()
		}
	}
}
