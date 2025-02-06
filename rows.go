package sqlslog

import (
	"database/sql/driver"
	"errors"
	"io"
	"log/slog"
	"reflect"
)

type rowsOptions struct {
	Close         StepOptions
	Next          StepOptions
	NextResultSet StepOptions
}

func defaultRowsOptions(formatter StepEventMsgBuilder) *rowsOptions {
	return &rowsOptions{
		Close:         *defaultStepOptions(formatter, StepRowsClose, LevelDebug),
		Next:          *defaultStepOptions(formatter, StepRowsNext, LevelDebug, HandleRowsNextError),
		NextResultSet: *defaultStepOptions(formatter, StepRowsNextResultSet, LevelDebug),
	}
}

func wrapRows(original driver.Rows, logger *stepLogger, options *rowsOptions) driver.Rows {
	if original == nil {
		return nil
	}
	rw := rowsWrapper{original: original, logger: logger, options: options}
	if rnrs, ok := original.(driver.RowsNextResultSet); ok {
		return &rowsNextResultSetWrapper{rw, rnrs}
	}
	return &rw
}

type rowsWrapper struct {
	original driver.Rows
	logger   *stepLogger
	options  *rowsOptions
}

var _ driver.Rows = (*rowsWrapper)(nil)

// Close implements driver.Rows.
func (r *rowsWrapper) Close() error {
	return ignoreAttr(r.logger.StepWithoutContext(&r.options.Close, withNilAttr(r.original.Close)))
}

// Columns implements driver.Rows.
func (r *rowsWrapper) Columns() []string {
	return r.original.Columns()
}

// Next implements driver.Rows.
func (r *rowsWrapper) Next(dest []driver.Value) error {
	return ignoreAttr(r.logger.StepWithoutContext(&r.options.Next, func() (*slog.Attr, error) {
		return nil, r.original.Next(dest)
	}))
}

// If the driver knows how to describe the types
// present in the returned result, it should implement the following
// interfaces: RowsColumnTypeScanType, RowsColumnTypeDatabaseTypeName,
// RowsColumnTypeLength, RowsColumnTypeNullable, and
// RowsColumnTypePrecisionScale. A given row value may also return a
// Rows type, which may represent a database cursor value.
//
// These are used in database/sql/sql.go
// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=3284-3300

var (
	_ driver.RowsColumnTypeScanType         = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeDatabaseTypeName = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeLength           = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypeNullable         = (*rowsWrapper)(nil)
	_ driver.RowsColumnTypePrecisionScale   = (*rowsWrapper)(nil)
)

// ColumnTypeScanType implements driver.RowsColumnTypeScanType.
func (r *rowsWrapper) ColumnTypeScanType(index int) reflect.Type {
	// https://cs.opensource.google/go/go/+/master:src/database/sql/sql.go;l=3284-3288
	if c, ok := r.original.(driver.RowsColumnTypeScanType); ok {
		return c.ColumnTypeScanType(index)
	}
	return reflect.TypeFor[any]()
}

// ColumnTypeDatabaseTypeName implements driver.RowsColumnTypeDatabaseTypeName.
func (r *rowsWrapper) ColumnTypeDatabaseTypeName(index int) string {
	if c, ok := r.original.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return c.ColumnTypeDatabaseTypeName(index)
	}
	return ""
}

// ColumnTypeLength implements driver.RowsColumnTypeLength.
func (r *rowsWrapper) ColumnTypeLength(index int) (int64, bool) {
	if c, ok := r.original.(driver.RowsColumnTypeLength); ok {
		return c.ColumnTypeLength(index)
	}
	return 0, false
}

// ColumnTypeNullable implements driver.RowsColumnTypeNullable.
func (r *rowsWrapper) ColumnTypeNullable(index int) (bool, bool) {
	if c, ok := r.original.(driver.RowsColumnTypeNullable); ok {
		return c.ColumnTypeNullable(index)
	}
	return false, false
}

// ColumnTypePrecisionScale implements driver.RowsColumnTypePrecisionScale.
func (r *rowsWrapper) ColumnTypePrecisionScale(index int) (int64, int64, bool) {
	if c, ok := r.original.(driver.RowsColumnTypePrecisionScale); ok {
		return c.ColumnTypePrecisionScale(index)
	}
	return 0, 0, false
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
	return ignoreAttr(
		r.logger.StepWithoutContext(
			&r.options.NextResultSet,
			withNilAttr(r.original.NextResultSet),
		),
	)
}

// HandleRowsNextError returns a boolean indicating completion and a slice of slog.Attr.
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
