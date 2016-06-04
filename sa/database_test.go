package sa

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/letsencrypt/boulder/test"
	"github.com/letsencrypt/boulder/test/vars"
)

func TestInvalidDSN(t *testing.T) {
	_, err := NewDbMap("invalid", 0)
	test.AssertError(t, err, "DB connect string missing the slash separating the database name")
}

var errExpected = errors.New("expected")

func TestMaxOpenConns(t *testing.T) {
	oldSetMaxOpenConns := setMaxOpenConns
	defer func() {
		setMaxOpenConns = oldSetMaxOpenConns
	}()
	maxOpenConns := -1
	setMaxOpenConns = func(db *sql.DB, m int) {
		maxOpenConns = m
		oldSetMaxOpenConns(db, maxOpenConns)
	}
	_, err := NewDbMap("sa@tcp(boulder-mysql:3306)/boulder_sa_integration", 100)
	if err != nil {
		t.Errorf("connecting to DB: %s", err)
	}
	if maxOpenConns != 100 {
		t.Errorf("maxOpenConns was not set: expected %d, got %d", 100, maxOpenConns)
	}
}

func TestNewDbMap(t *testing.T) {
	const mysqlConnectURL = "mysql+tcp://policy:password@boulder-mysql:3306/boulder_policy_integration?readTimeout=800ms&writeTimeout=800ms"
	const expectedTransformed = "policy:password@tcp(boulder-mysql:3306)/boulder_policy_integration?clientFoundRows=true&parseTime=true&readTimeout=800ms&strict=true&writeTimeout=800ms"

	oldSQLOpen := sqlOpen
	defer func() {
		sqlOpen = oldSQLOpen
	}()
	sqlOpen = func(dbType, connectString string) (*sql.DB, error) {
		if connectString != expectedTransformed {
			t.Errorf("incorrect connection string mangling, got %v", connectString)
		}
		return nil, errExpected
	}

	dbMap, err := NewDbMap(mysqlConnectURL, 0)
	if err != errExpected {
		t.Errorf("got incorrect error: %v", err)
	}
	if dbMap != nil {
		t.Errorf("expected nil, got %v", dbMap)
	}

}

func TestStrictness(t *testing.T) {
	const mysqlConnectURL = "mysql+tcp://policy:password@boulder-mysql:3306/boulder_policy_integration?readTimeout=800ms&writeTimeout=800ms"

	dbMap, err := NewDbMap(vars.DBConnSA, 1)
	_, err = dbMap.Exec(`insert into authz set
		id="hi", identifier="foo", status="pending", combinations="combos",
		registrationID=999999999999999999999999999;`)
	if err == nil {
		t.Fatal("Expected error when providing out of range value, got none.")
	}
	if !strings.Contains(err.Error(), "Out of range value for column") {
		t.Fatal("Got wrong type of error: %s", err)
	}
}
