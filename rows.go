package sqlslog

import (
	"database/sql/driver"
	"errors"
	"io"
	"log/slog"
	"reflect"
)

func wrapRows(original driver.Rows, logger *logger) driver.Rows {
	if original == nil {
		return nil
	}
	rw := rowsWrapper{original, logger}
	if rnrs, ok := original.(driver.RowsNextResultSet); ok {
		return &rowsNextResultSetWrapper{rw, rnrs}
	} else {
		return &rw
	}
}

type rowsWrapper struct {
	original driver.Rows
	logger   *logger
}

var _ driver.Rows = (*rowsWrapper)(nil)

// Close implements driver.Rows.
func (r *rowsWrapper) Close() error {
	return r.logger.StepWithoutContext(&r.logger.options.rowsClose, withNilAttr(r.original.Close))
}

// Columns implements driver.Rows.
func (r *rowsWrapper) Columns() []string {
	return r.original.Columns()
}

// Next implements driver.Rows.
func (r *rowsWrapper) Next(dest []driver.Value) error {
	return r.logger.StepWithoutContext(&r.logger.options.rowsNext, func() (*slog.Attr, error) {
		return nil, r.original.Next(dest)
	})
}

// If the driver knows how to describe the types
// present in the returned result it should implement the following
// interfaces: RowsColumnTypeScanType, RowsColumnTypeDatabaseTypeName,
// RowsColumnTypeLength, RowsColumnTypeNullable, and
// RowsColumnTypePrecisionScale. A given row value may also return a
// Rows type, which may represent a database cursor value.
//
// These are used in database/sql/sql.go
// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=3284-3300

var _ driver.RowsColumnTypeScanType = (*rowsWrapper)(nil)
var _ driver.RowsColumnTypeDatabaseTypeName = (*rowsWrapper)(nil)
var _ driver.RowsColumnTypeLength = (*rowsWrapper)(nil)
var _ driver.RowsColumnTypeNullable = (*rowsWrapper)(nil)
var _ driver.RowsColumnTypePrecisionScale = (*rowsWrapper)(nil)

// ColumnTypeScanType implements driver.RowsColumnTypeScanType.
func (r *rowsWrapper) ColumnTypeScanType(index int) reflect.Type {
	// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=3284-3288
	if c, ok := r.original.(driver.RowsColumnTypeScanType); ok {
		return c.ColumnTypeScanType(index)
	} else {
		return reflect.TypeFor[any]()
	}
}

// ColumnTypeDatabaseTypeName implements driver.RowsColumnTypeDatabaseTypeName.
func (r *rowsWrapper) ColumnTypeDatabaseTypeName(index int) string {
	if c, ok := r.original.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return c.ColumnTypeDatabaseTypeName(index)
	} else {
		return ""
	}
}

// ColumnTypeLength implements driver.RowsColumnTypeLength.
func (r *rowsWrapper) ColumnTypeLength(index int) (length int64, ok bool) {
	if c, ok := r.original.(driver.RowsColumnTypeLength); ok {
		return c.ColumnTypeLength(index)
	} else {
		return 0, false
	}
}

// ColumnTypeNullable implements driver.RowsColumnTypeNullable.
func (r *rowsWrapper) ColumnTypeNullable(index int) (nullable bool, ok bool) {
	if c, ok := r.original.(driver.RowsColumnTypeNullable); ok {
		return c.ColumnTypeNullable(index)
	} else {
		return false, false
	}
}

// ColumnTypePrecisionScale implements driver.RowsColumnTypePrecisionScale.
func (r *rowsWrapper) ColumnTypePrecisionScale(index int) (precision int64, scale int64, ok bool) {
	if c, ok := r.original.(driver.RowsColumnTypePrecisionScale); ok {
		return c.ColumnTypePrecisionScale(index)
	} else {
		return 0, 0, false
	}
}

type rowsNextResultSetWrapper struct {
	rowsWrapper
	original driver.RowsNextResultSet
}

// If multiple result sets are supported, Rows should implement
// RowsNextResultSet.
var _ driver.RowsNextResultSet = (*rowsNextResultSetWrapper)(nil)

// HasNextResultSet implements driver.RowsNextResultSet.
func (r *rowsNextResultSetWrapper) HasNextResultSet() bool {
	return r.original.HasNextResultSet()
}

// NextResultSet implements driver.RowsNextResultSet.
func (r *rowsNextResultSetWrapper) NextResultSet() error {
	return r.logger.StepWithoutContext(&r.logger.options.rowsNextResultSet, withNilAttr(r.original.NextResultSet))
}

// HandleRowsNextError returns completed and slice of slog.Attr.
// If err is nil, it returns true and a slice of slog.Attr{slog.Bool("eof", false)}.
// If err is io.EOF, it returns true and a slice of slog.Attr{slog.Bool("eof", true)}.
// Otherwise, it returns false and nil.
func HandleRowsNextError(err error) (bool, []slog.Attr) {
	if err == nil {
		return true, []slog.Attr{slog.Bool("eof", false)}
	}
	if errors.Is(err, io.EOF) {
		return true, []slog.Attr{slog.Bool("eof", true)}
	}
	return false, nil
}
