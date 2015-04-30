// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sa

import (
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/letsencrypt/boulder/core"
	blog "github.com/letsencrypt/boulder/log"
)

type SQLStorageAuthority struct {
	db     *sql.DB
	bucket map[string]interface{} // XXX included only for backward compat
	log    *blog.AuditLogger
}

func digest256(data []byte) []byte {
	d := sha256.New()
	_, _ = d.Write(data) // Never returns an error
	return d.Sum(nil)
}

func NewSQLStorageAuthority(logger *blog.AuditLogger, driver string, name string) (ssa *SQLStorageAuthority, err error) {
	logger.Notice("Storage Authority Starting")

	db, err := sql.Open(driver, name)
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}

	ssa = &SQLStorageAuthority{
		db:     db,
		log:    logger,
		bucket: make(map[string]interface{}),
	}

	err = ssa.InitTables()
	if err != nil {
		return
	}

	if driver == "mysql" {
		err = ssa.initIndexes()
		if err != nil {
			return
		}
	}

	return
}

func (ssa *SQLStorageAuthority) InitTables() (err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	// All fields should be created with "NOT NULL" because go's SQL support does not
	// handle null values well (see, e.g. https://github.com/go-sql-driver/mysql/issues/59)
	statements := []string{

	// Create registrations table
	`CREATE TABLE IF NOT EXISTS registrations (
		id VARCHAR(255) NOT NULL,
		thumbprint VARCHAR(255) NOT NULL,
		value BLOB NOT NULL
	);`,

	// Create pending authorizations table
	// TODO: Add NOT NULL to value. Right now it causes test failures because some
	// inserts to not fill all fields.
	`CREATE TABLE IF NOT EXISTS pending_authz (
		id VARCHAR(255) NOT NULL,
		value BLOB
	);`,

	// Create finalized authorizations table
	`CREATE TABLE IF NOT EXISTS authz (
		sequence INTEGER NOT NULL,
		id VARCHAR(255) NOT NULL,
		digest TEXT NOT NULL,
		value BLOB NOT NULL
	);`,

	// Create certificates table. This should be effectively append-only, enforced
	// by DB permissions.
	`CREATE TABLE IF NOT EXISTS certificates (
		serial VARCHAR(255) NOT NULL,
		digest VARCHAR(255) NOT NULL,
		value BLOB NOT NULL,
		issued DATETIME NOT NULL
		);`,

	// Create certificate status table. This provides metadata about a certificate
	// that can change over its lifetime, and rows are updateable unlike the
	// certificates table. The serial number primary key matches up with the one
	// on certificates.
	// subscriberApproved: 1 iff the subscriber has posted back to the server
	//   that they accept the certificate, otherwise 0.
	// status: 'good' or 'revoked'. Note that good, expired certificates remain
	//   with status 'good' but don't necessarily get fresh OCSP responses.
	// revokedDate: If status is 'revoked', this is the date and time it was
	//   revoked. Otherwise it has the zero value of time.Time, i.e. Jan 1 1970.
	// ocspLastUpdated: The date and time of the last time we generated an OCSP
	//   response. If we have never generated one, this has the zero value of
	//   time.Time, i.e. Jan 1 1970.
	`CREATE TABLE IF NOT EXISTS certificateStatus (
		serial VARCHAR(255) NOT NULL,
		subscriberApproved INTEGER NOT NULL,
		status VARCHAR(255) NOT NULL,
		revokedDate DATETIME NOT NULL,
		ocspLastUpdated DATETIME NOT NULL
		);`,
	}

	for _, statement := range statements {
		_, err = tx.Exec(statement)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	return
}

type indexTemplate struct {
	IndexName string
	TableName string
	Column    string
}

const sqlIndexTemplate = `SELECT IF (
    EXISTS(
        SELECT DISTINCT index_name FROM information_schema.statistics
        WHERE table_schema = DATABASE()
        AND table_name = '%s' AND index_name LIKE '%s'
    ),
    'SELECT 1;',
    'CREATE UNIQUE INDEX %s ON %s(%s)'
) INTO @a`

func (ssa *SQLStorageAuthority) initIndexes() (err error) {
	indexes := []indexTemplate{
		indexTemplate{
			IndexName: "registration_id",
			TableName: "registrations",
			Column:    "id",
		},
		indexTemplate{
			IndexName: "pending_authz_id",
			TableName: "pending_authz",
			Column:    "id",
		},
		indexTemplate{
			IndexName: "authz_id",
			TableName: "authz",
			Column:    "id",
		},
		indexTemplate{
			IndexName: "certificates_serial",
			TableName: "certificates",
			Column:    "serial",
		},
	}

	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	for _, template := range indexes {
		statement := fmt.Sprintf(sqlIndexTemplate, template.TableName, template.IndexName, template.IndexName, template.TableName, template.Column)
		_, err = tx.Exec(statement)
		if err != nil {
			tx.Rollback()
			return
		}
		_, err = tx.Exec("PREPARE ind_stmt FROM @a")
		if err != nil {
			tx.Rollback()
			return
		}
		_, err = tx.Exec("EXECUTE ind_stmt")
		if err != nil {
			tx.Rollback()
			return
		}
		_, err = tx.Exec("DEALLOCATE PREPARE ind_stmt")
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	return
}

func (ssa *SQLStorageAuthority) dumpTables(tx *sql.Tx) {
	fmt.Printf("===== TABLE DUMP =====\n")
	fmt.Printf("\n----- registrations -----\n")
	rows, err := tx.Query("SELECT id, thumbprint, value FROM registrations")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id, key, value []byte
			if err := rows.Scan(&id, &key, &value); err == nil {
				fmt.Printf("%s | %s | %s\n", string(id), string(key), hex.EncodeToString(value))
			} else {
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}

	fmt.Printf("\n----- pending_authz -----\n") // TODO
	fmt.Printf("\n----- authz -----\n")         // TODO
	fmt.Printf("\n----- certificates -----\n")  // TODO
}

func statusIsPending(status core.AcmeStatus) bool {
	return status == core.StatusPending || status == core.StatusProcessing || status == core.StatusUnknown
}

func existingPending(tx *sql.Tx, id string) (count int64) {
	tx.QueryRow("SELECT count(*) FROM pending_authz WHERE id = ?;", id).Scan(&count)
	return
}

func existingFinal(tx *sql.Tx, id string) (count int64) {
	tx.QueryRow("SELECT count(*) FROM authz WHERE id = ?;", id).Scan(&count)
	return
}

func existingRegistration(tx *sql.Tx, id string) (count int64) {
	tx.QueryRow("SELECT count(*) FROM registrations WHERE id = ?;", id).Scan(&count)
	return
}

func (ssa *SQLStorageAuthority) GetRegistration(id string) (reg core.Registration, err error) {
	var jsonReg []byte
	err = ssa.db.QueryRow("SELECT value FROM registrations WHERE id = ?;", id).Scan(&jsonReg)
	if err != nil {
		return
	}

	err = json.Unmarshal(jsonReg, &reg)
	return
}

func (ssa *SQLStorageAuthority) GetAuthorization(id string) (authz core.Authorization, err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	var jsonAuthz []byte
	err = tx.QueryRow("SELECT value FROM pending_authz WHERE id = ?;", id).Scan(&jsonAuthz)
	switch {
	case err == sql.ErrNoRows:
		err = tx.QueryRow("SELECT value FROM authz WHERE id = ?;", id).Scan(&jsonAuthz)
		if err != nil {
			tx.Rollback()
			return
		}
	case err != nil:
		tx.Rollback()
		return
	}
	tx.Commit()

	err = json.Unmarshal(jsonAuthz, &authz)
	return
}

// GetCertificate takes an id consisting of the first, sequential half of a
// serial number and returns the first certificate whose full serial number is
// lexically greater than that id. This allows clients to query on the known
// sequential half of our serial numbers to enumerate all certificates.
// TODO: Add index on certificates table
// TODO: Implement error when there are multiple certificates with the same
// sequential half.
func (ssa *SQLStorageAuthority) GetCertificate(id string) (cert []byte, err error) {
	if len(id) != 16 {
		err = errors.New("Invalid certificate serial " + id)
		return
	}
	err = ssa.db.QueryRow(
		"SELECT value FROM certificates WHERE serial LIKE ? LIMIT 1;",
		id + "%").Scan(&cert)
	return
}

// GetCertificateStatus takes a hexadecimal string representing the full 128-bit serial
// number of a certificate and returns data about that certificate's current
// validity.
func (ssa *SQLStorageAuthority) GetCertificateStatus(id string) (status core.CertificateStatus, err error) {
	if len(id) != 32 {
		err = errors.New("Invalid certificate serial " + id)
		return
	}
	var statusString string;
	err = ssa.db.QueryRow(
		`SELECT subscriberApproved, status, ocspLastUpdated
		 FROM certificateStatus
		 WHERE serial = ?
		 LIMIT 1;`, id).Scan(&status.SubscriberApproved, &statusString, &status.OCSPLastUpdated)
	if err != nil {
		return
	}
	status.Status = core.OCSPStatus(statusString)
	return
}

func (ssa *SQLStorageAuthority) NewRegistration() (id string, err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	// Check that it doesn't exist already
	candidate := core.NewToken()
	for existingRegistration(tx, candidate) > 0 {
		candidate = core.NewToken()
	}

	// Insert a stub row in pending
	_, err = tx.Exec("INSERT INTO registrations (id) VALUES (?);", candidate)
	if err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	id = candidate
	return
}

func (ssa *SQLStorageAuthority) UpdateRegistration(reg core.Registration) (err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	if existingRegistration(tx, reg.ID) != 1 {
		err = errors.New("Requested registration not found " + reg.ID)
		tx.Rollback()
		return
	}

	jsonReg, err := json.Marshal(reg)
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE registrations SET thumbprint=?, value=? WHERE id = ?;", reg.Key.Thumbprint, string(jsonReg), reg.ID)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}

func (ssa *SQLStorageAuthority) NewPendingAuthorization() (id string, err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	// Check that it doesn't exist already
	candidate := core.NewToken()
	for existingPending(tx, candidate) > 0 || existingFinal(tx, candidate) > 0 {
		candidate = core.NewToken()
	}

	// Insert a stub row in pending
	_, err = tx.Exec("INSERT INTO pending_authz (id) VALUES (?);", candidate)
	if err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	id = candidate
	return
}

func (ssa *SQLStorageAuthority) UpdatePendingAuthorization(authz core.Authorization) (err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	if !statusIsPending(authz.Status) {
		err = errors.New("Use Finalize() to update to a final status")
		tx.Rollback()
		return
	}

	if existingFinal(tx, authz.ID) > 0 {
		err = errors.New("Cannot update a final authorization")
		tx.Rollback()
		return
	}

	if existingPending(tx, authz.ID) != 1 {
		err = errors.New("Requested authorization not found " + authz.ID)
		tx.Rollback()
		return
	}

	jsonAuthz, err := json.Marshal(authz)
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec("UPDATE pending_authz SET value = ? WHERE id = ?;", jsonAuthz, authz.ID)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}

func (ssa *SQLStorageAuthority) FinalizeAuthorization(authz core.Authorization) (err error) {
	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	// Check that a pending authz exists
	if existingPending(tx, authz.ID) != 1 {
		err = errors.New("Cannot finalize a authorization that is not pending")
		tx.Rollback()
		return
	}
	if statusIsPending(authz.Status) {
		err = errors.New("Cannot finalize to a non-final status")
		tx.Rollback()
		return
	}

	// Manually set the index, to avoid AUTOINCREMENT issues
	var sequence int64
	var scanTarget sql.NullInt64
	err = tx.QueryRow("SELECT max(sequence) FROM authz").Scan(&scanTarget)
	switch {
	case !scanTarget.Valid:
		sequence = 0
	case err != nil:
		tx.Rollback()
		return
	default:
		sequence += scanTarget.Int64 + 1
	}

	jsonAuthz, err := json.Marshal(authz)
	if err != nil {
		tx.Rollback()
		return
	}
	digest := core.Fingerprint256(jsonAuthz)

	// Add to final table and delete from pending
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = tx.Exec("INSERT INTO authz (sequence, id, digest, value) VALUES (?, ?, ?, ?);", sequence, authz.ID, digest, jsonAuthz)
	if err != nil {
		tx.Rollback()
		return
	}
	_, err = tx.Exec("DELETE FROM pending_authz WHERE id = ?;", authz.ID)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

func (ssa *SQLStorageAuthority) AddCertificate(certDER []byte) (digest string, err error) {
	var parsedCertificate *x509.Certificate
	parsedCertificate, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}
	serial := fmt.Sprintf("%032x", parsedCertificate.SerialNumber)

	tx, err := ssa.db.Begin()
	if err != nil {
		return
	}

	digest = core.Fingerprint256(certDER)
	_, err = tx.Exec("INSERT INTO certificates (serial, digest, value, issued) VALUES (?,?,?,?);",
		serial, digest, certDER, time.Now())
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = tx.Exec(`
		INSERT INTO certificateStatus
		(serial, subscriberApproved, status, revokedDate, ocspLastUpdated)
		VALUES (?, 0, 'good', ?, ?);
		`, serial, time.Time{}, time.Time{})
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	return
}
