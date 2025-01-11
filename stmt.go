package sqlslog

import (
	"context"
	"database/sql/driver"
	"log/slog"
)

type stmtWrapper struct {
	original driver.Stmt
	logger   *slog.Logger
}

var _ driver.Stmt = (*stmtWrapper)(nil)

// Close implements driver.Stmt.
func (s *stmtWrapper) Close() error {
	panic("unimplemented")
}

// Exec implements driver.Stmt.
func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	panic("unimplemented")
}

// NumInput implements driver.Stmt.
func (s *stmtWrapper) NumInput() int {
	panic("unimplemented")
}

// Query implements driver.Stmt.
func (s *stmtWrapper) Query(args []driver.Value) (driver.Rows, error) {
	panic("unimplemented")
}

type stmtExecContextWrapperImpl struct {
	original driver.StmtExecContext
	logger   *slog.Logger
}

var _ driver.StmtExecContext = (*stmtExecContextWrapperImpl)(nil)

// ExecContext implements driver.StmtExecContext.
func (s *stmtExecContextWrapperImpl) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	panic("unimplemented")
}

type stmtQueryContextWrapperImpl struct {
	original driver.StmtQueryContext
	logger   *slog.Logger
}

var _ driver.StmtQueryContext = (*stmtQueryContextWrapperImpl)(nil)

// QueryContext implements driver.StmtQueryContext.
func (s *stmtQueryContextWrapperImpl) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	panic("unimplemented")
}

type stmtExecContextWrapper struct {
	stmtWrapper
	stmtExecContextWrapperImpl
}

var _ driver.Stmt = (*stmtExecContextWrapper)(nil)
var _ driver.StmtExecContext = (*stmtExecContextWrapper)(nil)

type stmtQueryContextWrapper struct {
	stmtWrapper
	stmtQueryContextWrapperImpl
}

var _ driver.Stmt = (*stmtQueryContextWrapper)(nil)
var _ driver.StmtQueryContext = (*stmtQueryContextWrapper)(nil)

type stmtContextWrapper struct {
	stmtWrapper
	stmtExecContextWrapperImpl
	stmtQueryContextWrapperImpl
}

var _ driver.Stmt = (*stmtContextWrapper)(nil)
var _ driver.StmtExecContext = (*stmtContextWrapper)(nil)
var _ driver.StmtQueryContext = (*stmtContextWrapper)(nil)

func wrapStmt(original driver.Stmt, logger *slog.Logger) driver.Stmt {
	if original == nil {
		return nil
	}
	stmtWrapper := stmtWrapper{original: original, logger: logger}

	stmtExec, withExecContext := original.(driver.StmtExecContext)
	stmtQuery, withQueryContext := original.(driver.StmtQueryContext)
	if withExecContext && withQueryContext {
		return &stmtContextWrapper{
			stmtWrapper:                 stmtWrapper,
			stmtExecContextWrapperImpl:  stmtExecContextWrapperImpl{original: stmtExec, logger: logger},
			stmtQueryContextWrapperImpl: stmtQueryContextWrapperImpl{original: stmtQuery, logger: logger},
		}
	}
	if withExecContext {
		return &stmtExecContextWrapper{
			stmtWrapper:                stmtWrapper,
			stmtExecContextWrapperImpl: stmtExecContextWrapperImpl{original: stmtExec, logger: logger},
		}
	}
	if withQueryContext {
		return &stmtQueryContextWrapper{
			stmtWrapper:                 stmtWrapper,
			stmtQueryContextWrapperImpl: stmtQueryContextWrapperImpl{original: stmtQuery, logger: logger},
		}
	}
	return &stmtWrapper
}
