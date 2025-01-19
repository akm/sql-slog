package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

type stmtWrapper struct {
	original driver.Stmt
	logger   *logger
}

var _ driver.Stmt = (*stmtWrapper)(nil)

// Close implements driver.Stmt.
func (s *stmtWrapper) Close() error {
	return s.logger.StepWithoutContext(&s.logger.options.stmtClose, s.original.Close)
}

// Exec implements driver.Stmt.
func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := lg.StepWithoutContext(&s.logger.options.stmtExec, func() error {
		var err error
		result, err = s.original.Exec(args) //nolint:staticcheck
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// NumInput implements driver.Stmt.
func (s *stmtWrapper) NumInput() int {
	return s.original.NumInput()
}

// Query implements driver.Stmt.
func (s *stmtWrapper) Query(args []driver.Value) (driver.Rows, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var rows driver.Rows
	err := lg.StepWithoutContext(&s.logger.options.stmtQuery, func() error {
		var err error
		rows, err = s.original.Query(args) //nolint:staticcheck
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger), nil
}

type stmtExecContextWrapperImpl struct {
	original driver.StmtExecContext
	logger   *logger
}

var _ driver.StmtExecContext = (*stmtExecContextWrapperImpl)(nil)

// ExecContext implements driver.StmtExecContext.
func (s *stmtExecContextWrapperImpl) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := lg.logActionContext(ctx, &s.logger.options.stmtExecContext, func() error {
		var err error
		result, err = s.original.ExecContext(ctx, args)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

type stmtQueryContextWrapperImpl struct {
	original driver.StmtQueryContext
	logger   *logger
}

var _ driver.StmtQueryContext = (*stmtQueryContextWrapperImpl)(nil)

// QueryContext implements driver.StmtQueryContext.
func (s *stmtQueryContextWrapperImpl) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var rows driver.Rows
	err := lg.logActionContext(ctx, &s.logger.options.stmtQueryContext, func() error {
		var err error
		rows, err = s.original.QueryContext(ctx, args)
		return err
	})
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger), nil
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

func wrapStmt(original driver.Stmt, logger *logger) driver.Stmt {
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
	// Commented out because the original implementation does not have this check.
	//
	// if withExecContext {
	// 	return &stmtExecContextWrapper{
	// 		stmtWrapper:                stmtWrapper,
	// 		stmtExecContextWrapperImpl: stmtExecContextWrapperImpl{original: stmtExec, logger: logger},
	// 	}
	// }
	// if withQueryContext {
	// 	return &stmtQueryContextWrapper{
	// 		stmtWrapper:                 stmtWrapper,
	// 		stmtQueryContextWrapperImpl: stmtQueryContextWrapperImpl{original: stmtQuery, logger: logger},
	// 	}
	// }
	return &stmtWrapper
}
