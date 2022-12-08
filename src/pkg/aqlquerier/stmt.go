package aqlquerier

import (
	"context"
	"database/sql/driver"
	"errors"
	"hms/gateway/pkg/aqlprocessor"
	"log"
)

type Stmt struct {
	query *aqlprocessor.Query
}

// Close closes the statement.
// As of Go 1.1, a Stmt will not be closed if it's in use
// by any queries.
//
// Drivers must ensure all network calls made by Close
// do not block indefinitely (e.g. apply a timeout).
func (stmt *Stmt) Close() error {
	return nil
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (stmt *Stmt) NumInput() int {
	return stmt.query.ParametersCount() // TODO: Implement
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
//
// Deprecated: Drivers should implement StmtExecContext instead (or additionally).
func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("Query.Exec not implemented") // nolint
}

// Query executes a query that may return rows, such as a
// SELECT.
//
// Deprecated: Drivers should implement StmtQueryContext instead (or additionally).
func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("Stmt.Query deprecated and not implemented")
}

// QueryContext executes a query that may return rows, such as a
// SELECT.
//
// QueryContext must honor the context timeout and return when it is canceled.
func (stmt *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	for _, arg := range args {
		log.Printf("arg: %v = %v (%v)", arg.Name, arg.Value, arg.Ordinal)
	}

	rows := Rows{}

	for _, se := range stmt.query.Select.SelectExprs {
		rows.columns = append(rows.columns, se.AliasName)
	}

	rows.rows = append(rows.rows,
		Row{values: []interface{}{123}},
		Row{values: []interface{}{256}},
	)
	return &rows, nil
}
