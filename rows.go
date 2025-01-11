package sqlslog

import (
	"database/sql/driver"
	"reflect"
)

type rows struct {
}

var _ driver.Rows = (*rows)(nil)

// If multiple result sets are supported, Rows should implement
// RowsNextResultSet. If the driver knows how to describe the types
// present in the returned result it should implement the following
// interfaces: RowsColumnTypeScanType, RowsColumnTypeDatabaseTypeName,
// RowsColumnTypeLength, RowsColumnTypeNullable, and
// RowsColumnTypePrecisionScale. A given row value may also return a
// Rows type, which may represent a database cursor value.

var _ driver.RowsNextResultSet = (*rows)(nil)
var _ driver.RowsColumnTypeScanType = (*rows)(nil)
var _ driver.RowsColumnTypeDatabaseTypeName = (*rows)(nil)
var _ driver.RowsColumnTypeLength = (*rows)(nil)
var _ driver.RowsColumnTypeNullable = (*rows)(nil)
var _ driver.RowsColumnTypePrecisionScale = (*rows)(nil)

// ColumnTypePrecisionScale implements driver.RowsColumnTypePrecisionScale.
func (r *rows) ColumnTypePrecisionScale(index int) (precision int64, scale int64, ok bool) {
	panic("unimplemented")
}

// ColumnTypeNullable implements driver.RowsColumnTypeNullable.
func (r *rows) ColumnTypeNullable(index int) (nullable bool, ok bool) {
	panic("unimplemented")
}

// ColumnTypeLength implements driver.RowsColumnTypeLength.
func (r *rows) ColumnTypeLength(index int) (length int64, ok bool) {
	panic("unimplemented")
}

// ColumnTypeDatabaseTypeName implements driver.RowsColumnTypeDatabaseTypeName.
func (r *rows) ColumnTypeDatabaseTypeName(index int) string {
	panic("unimplemented")
}

// ColumnTypeScanType implements driver.RowsColumnTypeScanType.
func (r *rows) ColumnTypeScanType(index int) reflect.Type {
	panic("unimplemented")
}

// Close implements driver.RowsNextResultSet.
func (r *rows) Close() error {
	panic("unimplemented")
}

// Columns implements driver.RowsNextResultSet.
func (r *rows) Columns() []string {
	panic("unimplemented")
}

// HasNextResultSet implements driver.RowsNextResultSet.
func (r *rows) HasNextResultSet() bool {
	panic("unimplemented")
}

// Next implements driver.RowsNextResultSet.
func (r *rows) Next(dest []driver.Value) error {
	panic("unimplemented")
}

// NextResultSet implements driver.RowsNextResultSet.
func (r *rows) NextResultSet() error {
	panic("unimplemented")
}
