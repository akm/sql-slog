package wrap

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

type stmtWrapper struct {
	original driver.Stmt
	logger   *SQLLogger
}

var _ driver.Stmt = (*stmtWrapper)(nil)

// Close implements driver.Stmt.
func (s *stmtWrapper) Close() error {
	return IgnoreAttr(s.logger.StepWithoutContext(&s.logger.Options.StmtClose, WithNilAttr(s.original.Close)))
}

// Exec implements driver.Stmt.
func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := IgnoreAttr(lg.StepWithoutContext(&s.logger.Options.StmtExec, func() (*slog.Attr, error) {
		var err error
		result, err = s.original.Exec(args) //nolint:staticcheck
		return nil, err
	}))
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
	err := IgnoreAttr(lg.StepWithoutContext(&s.logger.Options.StmtQuery, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.Query(args) //nolint:staticcheck
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger), nil
}

type stmtExecContextWrapperImpl struct {
	original driver.StmtExecContext
	logger   *SQLLogger
}

var _ driver.StmtExecContext = (*stmtExecContextWrapperImpl)(nil)

// ExecContext implements driver.StmtExecContext.
func (s *stmtExecContextWrapperImpl) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := IgnoreAttr(lg.Step(ctx, &s.logger.Options.StmtExecContext, func() (*slog.Attr, error) {
		var err error
		result, err = s.original.ExecContext(ctx, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return result, nil
}

type stmtQueryContextWrapperImpl struct {
	original driver.StmtQueryContext
	logger   *SQLLogger
}

var _ driver.StmtQueryContext = (*stmtQueryContextWrapperImpl)(nil)

// QueryContext implements driver.StmtQueryContext.
func (s *stmtQueryContextWrapperImpl) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var rows driver.Rows
	err := IgnoreAttr(lg.Step(ctx, &s.logger.Options.StmtQueryContext, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.QueryContext(ctx, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger), nil
}

type stmtExecContextWrapper struct {
	stmtWrapper
	stmtExecContextWrapperImpl
}

var (
	_ driver.Stmt            = (*stmtExecContextWrapper)(nil)
	_ driver.StmtExecContext = (*stmtExecContextWrapper)(nil)
)

type stmtQueryContextWrapper struct {
	stmtWrapper
	stmtQueryContextWrapperImpl
}

var (
	_ driver.Stmt             = (*stmtQueryContextWrapper)(nil)
	_ driver.StmtQueryContext = (*stmtQueryContextWrapper)(nil)
)

type stmtContextWrapper struct {
	stmtWrapper
	stmtExecContextWrapperImpl
	stmtQueryContextWrapperImpl
}

var (
	_ driver.Stmt             = (*stmtContextWrapper)(nil)
	_ driver.StmtExecContext  = (*stmtContextWrapper)(nil)
	_ driver.StmtQueryContext = (*stmtContextWrapper)(nil)
)

func wrapStmt(original driver.Stmt, logger *SQLLogger) driver.Stmt {
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
