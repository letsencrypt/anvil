package db

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	gorp "gopkg.in/go-gorp/gorp.v2"
)

// ErrDatabaseOp wraps an underlying err with a description of the operation
// that was being performed when the error occurred (insert, select, select
// one, exec, etc) and the table that the operation was being performed on.
type ErrDatabaseOp struct {
	Op    string
	Table string
	Err   error
}

// NoRows returns true when the underlying error is sql.ErrNoRows and indicates
// that the error was that no results were found.
func (e ErrDatabaseOp) NoRows() bool {
	return e.Err == sql.ErrNoRows
}

// Duplicate returns true when the underlying error has a message with a prefix
// matching "Error 1062: Duplicate entry". This is the error prefixed returned
// by MariaDB when a duplicate row is to be inserted.
func (e ErrDatabaseOp) Duplicate() bool {
	return strings.HasPrefix(
		e.Err.Error(),
		"Error 1062: Duplicate entry")
}

// Error for an ErrDatabaseOp composes a message with context about the
// operation and table as well as the underlying Err's error message.
func (e ErrDatabaseOp) Error() string {
	// If there is a table, include it in the context
	if e.Table != "" {
		return fmt.Sprintf(
			"failed to %s %s: %s",
			e.Op,
			e.Table,
			e.Err)
	}
	return fmt.Sprintf(
		"failed to %s: %s",
		e.Op,
		e.Err)
}

// IsNoRowsErr is a utility function for casting an error to ErrDatabaseOp and
// returning the result of its NoRows() function. If the error is not an
// ErrDatabaseOp the return value of IsNoRowsErr will always be false.
func IsNoRowsErr(err error) bool {
	// if the err is an ErrDatabaseOp instance, return its NoRows() result to see
	// if the inner err is sql.ErrNoRows
	if dbErr, ok := err.(ErrDatabaseOp); ok {
		return dbErr.NoRows()
	}
	return false
}

// IsDuplicateErr is a utility function for casting an error to ErrDatabaseOp and
// returning the result of its Duplicate() function. If the error is not an
// ErrDatabaseOp the return value of IsDuplicateErr will always be false.
func IsDuplicateErr(err error) bool {
	// if the err is an ErrDatabaseOp instance, return its Duplicate() result to
	// see if the inner err indicates a duplicate row error.
	if dbErr, ok := err.(ErrDatabaseOp); ok {
		return dbErr.Duplicate()
	}
	return false
}

// WrappedMap wraps a *gorp.DbMap such that its major functions wrap error
// results in ErrDatabaseOp instances before returning them to the caller.
type WrappedMap struct {
	*gorp.DbMap
}

func (m *WrappedMap) Get(holder interface{}, keys ...interface{}) (interface{}, error) {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Get(holder, keys...)
}

func (m *WrappedMap) Insert(list ...interface{}) error {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Insert(list...)
}

func (m *WrappedMap) Update(list ...interface{}) (int64, error) {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Update(list...)
}

func (m *WrappedMap) Delete(list ...interface{}) (int64, error) {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Delete(list...)
}

func (m *WrappedMap) Select(holder interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Select(holder, query, args...)
}

func (m *WrappedMap) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return WrappedExecutor{SqlExecutor: m.DbMap}.SelectOne(holder, query, args...)
}

func (m *WrappedMap) Exec(query string, args ...interface{}) (sql.Result, error) {
	return WrappedExecutor{SqlExecutor: m.DbMap}.Exec(query, args...)
}

func (m *WrappedMap) WithContext(ctx context.Context) gorp.SqlExecutor {
	return WrappedExecutor{SqlExecutor: m.DbMap.WithContext(ctx)}
}

func (m *WrappedMap) Begin() (Transaction, error) {
	tx, err := m.DbMap.Begin()
	if err != nil {
		return tx, ErrDatabaseOp{
			Op:  "begin transaction",
			Err: err,
		}
	}
	return WrappedTransaction{
		Transaction: tx,
	}, err
}

// WrappedTransaction wraps a *gorp.Transaction such that its major functions
// wrap error results in ErrDatabaseOp instances before returning them to the
// caller.
type WrappedTransaction struct {
	*gorp.Transaction
}

func (tx WrappedTransaction) WithContext(ctx context.Context) gorp.SqlExecutor {
	return WrappedExecutor{SqlExecutor: tx.Transaction.WithContext(ctx)}
}

func (tx WrappedTransaction) Commit() error {
	return tx.Transaction.Commit()
}

func (tx WrappedTransaction) Rollback() error {
	return tx.Transaction.Rollback()
}

func (tx WrappedTransaction) Get(holder interface{}, keys ...interface{}) (interface{}, error) {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Get(holder, keys...)
}

func (tx WrappedTransaction) Insert(list ...interface{}) error {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Insert(list...)
}

func (tx WrappedTransaction) Update(list ...interface{}) (int64, error) {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Update(list...)
}

func (tx WrappedTransaction) Delete(list ...interface{}) (int64, error) {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Delete(list...)
}

func (tx WrappedTransaction) Select(holder interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Select(holder, query, args...)
}

func (tx WrappedTransaction) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).SelectOne(holder, query, args...)
}

func (tx WrappedTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return (WrappedExecutor{SqlExecutor: tx.Transaction}).Exec(query, args...)
}

var (
	// selectTableRegexp matches the table name from an SQL select statement
	selectTableRegexp = regexp.MustCompile(`(?i)^\s*select\s+[a-z\d\.\(\), \_\*` + "`" + `]+\s+from\s+([a-z\d\_,` + "`" + `]+)`)
	// insertTableRegexp matches the table name from an SQL insert statement
	insertTableRegexp = regexp.MustCompile(`(?i)^\s*insert\s+into\s+([a-z\d \_,` + "`" + `]+)\s+(?:set|\()`)
	// updateTableRegexp matches the table name from an SQL update statement
	updateTableRegexp = regexp.MustCompile(`(?i)^\s*update\s+([a-z\d \_,` + "`" + `]+)\s+set`)
	// deleteTableRegexp matches the table name from an SQL delete statement
	deleteTableRegexp = regexp.MustCompile(`(?i)^\s*delete\s+from\s+([a-z\d \_,` + "`" + `]+)\s+where`)

	// tableRegexps is a list of regexps that tableFromQuery will try to use in
	// succession to find the table name for an SQL query.
	tableRegexps = []*regexp.Regexp{
		selectTableRegexp,
		insertTableRegexp,
		updateTableRegexp,
		deleteTableRegexp,
	}
)

// tableFromQuery uses the tableRegexps on the provided query to return the
// associated table name or an empty string if it can't be determined from the
// query.
func tableFromQuery(query string) string {
	for _, r := range tableRegexps {
		if matches := r.FindStringSubmatch(query); len(matches) >= 2 {
			return matches[1]
		}
	}
	return ""
}

// WrappedExecutor wraps a gorp.SqlExecutor such that its major functions
// wrap error results in ErrDatabaseOp instances before returning them to the
// caller.
type WrappedExecutor struct {
	gorp.SqlExecutor
}

func (we WrappedExecutor) Get(holder interface{}, keys ...interface{}) (interface{}, error) {
	res, err := we.SqlExecutor.Get(holder, keys...)
	if err != nil {
		return res, ErrDatabaseOp{
			Op:    "get",
			Table: fmt.Sprintf("%T", holder),
			Err:   err,
		}
	}
	return res, err
}

func (we WrappedExecutor) Insert(list ...interface{}) error {
	if err := we.SqlExecutor.Insert(list...); err != nil {
		table := "unknown"
		if len(list) > 0 {
			table = fmt.Sprintf("%T", list[0])
		}
		return ErrDatabaseOp{
			Op:    "insert",
			Table: table,
			Err:   err,
		}
	}
	return nil
}

func (we WrappedExecutor) Update(list ...interface{}) (int64, error) {
	updatedRows, err := we.SqlExecutor.Update(list...)
	if err != nil {
		table := "unknown"
		if len(list) > 0 {
			table = fmt.Sprintf("%T", list[0])
		}
		return updatedRows, ErrDatabaseOp{
			Op:    "update",
			Table: table,
			Err:   err,
		}
	}
	return updatedRows, err
}

func (we WrappedExecutor) Delete(list ...interface{}) (int64, error) {
	deletedRows, err := we.SqlExecutor.Delete(list...)
	if err != nil {
		table := "unknown"
		if len(list) > 0 {
			table = fmt.Sprintf("%T", list[0])
		}
		return deletedRows, ErrDatabaseOp{
			Op:    "delete",
			Table: table,
			Err:   err,
		}
	}
	return deletedRows, err
}

func (we WrappedExecutor) Select(holder interface{}, query string, args ...interface{}) ([]interface{}, error) {
	result, err := we.SqlExecutor.Select(holder, query, args...)
	if err != nil {
		table := fmt.Sprintf("%T", holder)
		if extractedTable := tableFromQuery(query); extractedTable != "" {
			table = extractedTable
		}
		return result, ErrDatabaseOp{
			Op:    "select",
			Table: table,
			Err:   err,
		}
	}
	return result, err
}

func (we WrappedExecutor) SelectOne(holder interface{}, query string, args ...interface{}) error {
	if err := we.SqlExecutor.SelectOne(holder, query, args...); err != nil {
		table := fmt.Sprintf("%T", holder)
		if extractedTable := tableFromQuery(query); extractedTable != "" {
			table = extractedTable
		}
		return ErrDatabaseOp{
			Op:    "select one",
			Table: table,
			Err:   err,
		}
	}
	return nil
}

func (we WrappedExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := we.SqlExecutor.Exec(query, args...)
	if err != nil {
		table := "unknown"
		if extractedTable := tableFromQuery(query); extractedTable != "" {
			table = extractedTable
		}
		return nil, ErrDatabaseOp{
			Op:    "exec",
			Table: table,
			Err:   err,
		}
	}
	return res, nil
}