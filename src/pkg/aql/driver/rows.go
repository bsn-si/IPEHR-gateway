package driver

import (
	"database/sql/driver"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type AQLRowser interface {
	driver.Rows
	NamedColumns() []Column
}

type Column struct {
	Name string
	Path string
}

type Rows struct {
	rows    []Row
	columns []Column

	cursor int
}

type Row struct {
	values []interface{}
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice. If a particular column name isn't known, an empty
// string should be returned for that entry.
func (rs *Rows) Columns() []string {
	result := make([]string, 0, len(rs.columns))
	for _, c := range rs.columns {
		result = append(result, c.Name)
	}

	return result
}

// Close closes the rows iterator.
func (rs *Rows) Close() error {
	return nil
}

func (rs *Rows) NamedColumns() []Column {
	return rs.columns
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (rs *Rows) Next(dest []driver.Value) error {
	if len(rs.rows) <= rs.cursor {
		return errors.New("now rows")
	}

	row := rs.rows[rs.cursor]
	if len(dest) != len(row.values) {
		return errors.New("invalid count of values in row")
	}

	for i, v := range row.values {
		dest[i] = driver.Value(v)
	}

	rs.cursor++
	return nil
}
