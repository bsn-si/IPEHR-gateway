package driver

import (
	"context"
	"database/sql"
	"database/sql/driver"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"
)

func init() {
	sql.Register("aql", &AQLDriver{})
}

type AQLDriver struct{}

func (svc *AQLDriver) Open(name string) (driver.Conn, error) {
	conn := &AQLConn{
		index: treeindex.DefaultEHRIndex,
	}

	return conn, nil
}

type AQLConn struct {
	index *treeindex.EHRIndex
}

// Prepare returns a prepared statement, bound to this connection.
func (conn *AQLConn) Prepare(query string) (driver.Stmt, error) {
	return conn.PrepareContext(context.Background(), query)
}

func (conn *AQLConn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	aqlQuery, err := aqlprocessor.NewAqlProcessor(query).Process()
	if err != nil {
		return nil, errors.Wrap(err, "cannot prepare AQL query")
	}

	stmt := &Stmt{
		query: aqlQuery,
		index: conn.index,
	}

	return stmt, nil
}

// implements driver.Tx interface
// Rollback - not implemented ...
func (conn *AQLConn) Rollback() error {
	return errors.New("Rollback method not implemented")
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
//
// Drivers must ensure all network calls made by Close
// do not block indefinitely (e.g. apply a timeout).
func (conn *AQLConn) Close() error {
	return nil
}

// Begin starts and returns a new transaction.
//
// Deprecated: Drivers should implement ConnBeginTx instead (or additionally).
func (conn *AQLConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}
