package sa

import (
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/letsencrypt/boulder/core"
	corepb "github.com/letsencrypt/boulder/core/proto"
	berrors "github.com/letsencrypt/boulder/errors"
	"github.com/letsencrypt/boulder/features"
	bgrpc "github.com/letsencrypt/boulder/grpc"
	sapb "github.com/letsencrypt/boulder/sa/proto"
)

var errIncompleteRequest = errors.New("Incomplete gRPC request message")

// AddSerial writes a record of a serial number generation to the DB.
func (ssa *SQLStorageAuthority) AddSerial(ctx context.Context, req *sapb.AddSerialRequest) (*corepb.Empty, error) {
	if req == nil || req.Created == nil || req.Expires == nil || req.Serial == nil || req.RegID == nil {
		return nil, errIncompleteRequest
	}
	created := time.Unix(0, *req.Created)
	expires := time.Unix(0, *req.Expires)
	err := ssa.dbMap.WithContext(ctx).Insert(&recordedSerialModel{
		Serial:         *req.Serial,
		RegistrationID: *req.RegID,
		Created:        created,
		Expires:        expires,
	})
	if err != nil {
		return nil, err
	}
	return &corepb.Empty{}, nil
}

// AddPrecertificate writes a record of a precertificate generation to the DB.
func (ssa *SQLStorageAuthority) AddPrecertificate(ctx context.Context, req *sapb.AddCertificateRequest) (*corepb.Empty, error) {
	if req == nil || req.Der == nil || req.Issued == nil || req.RegID == nil {
		return nil, errIncompleteRequest
	}
	parsed, err := x509.ParseCertificate(req.Der)
	if err != nil {
		return nil, err
	}
	issued := time.Unix(0, *req.Issued)
	serialHex := core.SerialToString(parsed.SerialNumber)
	err = ssa.dbMap.WithContext(ctx).Insert(&precertificateModel{
		Serial:         serialHex,
		RegistrationID: *req.RegID,
		DER:            req.Der,
		Issued:         issued,
		Expires:        parsed.NotAfter,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			return nil, berrors.DuplicateError("cannot add a duplicate precertificate")
		}
		return nil, err
	}

	// With feature.StoreIssuerInfo we've added a new field to certStatusModel
	// so when we try and use dbMap.Insert it will always try to insert that field.
	// That will break when the relevant migration hasn't bene applied so we need
	// to use an explicit INSERT statement that we can manipulate to include the
	// field only when the feature is enabled (and as such the migration has been
	// applied).
	csFields := certStatusFields
	if features.Enabled(features.StoreIssuerInfo) {
		csFields += ", issuerID"
	}
	qmarks := []string{}
	for range strings.Split(csFields, ",") {
		qmarks = append(qmarks, "?")
	}
	args := []interface{}{
		serialHex,                   // serial
		string(core.OCSPStatusGood), // stauts
		ssa.clk.Now(),               // ocspLastUpdated
		time.Time{},                 // revokedDate
		0,                           // revokedReason
		time.Time{},                 // lastExpirationNagSent
		req.Ocsp,                    // ocspResponse
		parsed.NotAfter,             // notAfter
		false,                       // isExpired
	}
	if features.Enabled(features.StoreIssuerInfo) && req.IssuerID != nil {
		args = append(args, req.IssuerID)
	}

	_, err = ssa.dbMap.WithContext(ctx).Exec(fmt.Sprintf(
		"INSERT INTO certificateStatus (%s) VALUES (%s)",
		csFields,
		strings.Join(qmarks, ","),
	), args...)
	if err != nil {
		return nil, err
	}
	return &corepb.Empty{}, nil
}

// GetPrecertificate takes a serial number and returns the corresponding
// precertificate, or error if it does not exist.
func (ssa *SQLStorageAuthority) GetPrecertificate(ctx context.Context, reqSerial *sapb.Serial) (*corepb.Certificate, error) {
	if !core.ValidSerial(*reqSerial.Serial) {
		return nil,
			fmt.Errorf("Invalid precertificate serial %q", *reqSerial.Serial)
	}
	cert, err := SelectPrecertificate(ssa.dbMap.WithContext(ctx), *reqSerial.Serial)
	if err == sql.ErrNoRows {
		return nil,
			berrors.NotFoundError("precertificate with serial %q not found", *reqSerial.Serial)
	}
	if err != nil {
		return nil, err
	}

	return bgrpc.CertToPB(cert), nil
}
