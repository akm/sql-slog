package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

type stmtOptions struct {
	Close        StepOptions
	Exec         StepOptions
	Query        StepOptions
	ExecContext  StepOptions
	QueryContext StepOptions

	Rows *rowsOptions
}

func defaultStmtOptions(msgb StepEventMsgBuilder) *stmtOptions {
	return &stmtOptions{
		Close:        *defaultStepOptions(msgb, StepStmtClose, LevelInfo),
		Exec:         *defaultStepOptions(msgb, StepStmtExec, LevelInfo),
		Query:        *defaultStepOptions(msgb, StepStmtQuery, LevelInfo),
		ExecContext:  *defaultStepOptions(msgb, StepStmtExecContext, LevelInfo),
		QueryContext: *defaultStepOptions(msgb, StepStmtQueryContext, LevelInfo),
		Rows:         defaultRowsOptions(msgb),
	}
}

func wrapStmt(original driver.Stmt, logger *stepLogger, options *stmtOptions) driver.Stmt {
	if original == nil {
		return nil
	}
	stmtWrapper := stmtWrapper{original: original, logger: logger, options: options}

	stmtExec, withExecContext := original.(driver.StmtExecContext)
	stmtQuery, withQueryContext := original.(driver.StmtQueryContext)
	if withExecContext && withQueryContext {
		stmtCtxW := &stmtContextWrapper{
			stmtWrapper:                 stmtWrapper,
			stmtExecContextWrapperImpl:  stmtExecContextWrapperImpl{original: stmtExec, logger: logger, options: options},
			stmtQueryContextWrapperImpl: stmtQueryContextWrapperImpl{original: stmtQuery, logger: logger, options: options},
		}
		if nvc, ok := original.(driver.NamedValueChecker); ok {
			return &stmtContextNvcWrapper{
				stmtContextWrapper: *stmtCtxW,
				NamedValueChecker:  nvc,
			}
		}
		return stmtCtxW
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

	if nvc, ok := original.(driver.NamedValueChecker); ok {
		return &stmtNvcWrapper{
			stmtWrapper:       stmtWrapper,
			NamedValueChecker: nvc,
		}
	}
	return &stmtWrapper
}

type stmtWrapper struct {
	original driver.Stmt
	logger   *stepLogger
	options  *stmtOptions
}

var _ driver.Stmt = (*stmtWrapper)(nil)

// Close implements driver.Stmt.
func (s *stmtWrapper) Close() error {
	return ignoreAttr(s.logger.StepWithoutContext(&s.options.Close, withNilAttr(s.original.Close)))
}

// Exec implements driver.Stmt.
func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := ignoreAttr(lg.StepWithoutContext(&s.options.Exec, func() (*slog.Attr, error) {
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
	err := ignoreAttr(lg.StepWithoutContext(&s.options.Query, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.Query(args) //nolint:staticcheck
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger, s.options.Rows), nil
}

type stmtExecContextWrapperImpl struct {
	original driver.StmtExecContext
	logger   *stepLogger
	options  *stmtOptions
}

var _ driver.StmtExecContext = (*stmtExecContextWrapperImpl)(nil)

// ExecContext implements driver.StmtExecContext.
func (s *stmtExecContextWrapperImpl) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := ignoreAttr(lg.Step(ctx, &s.options.ExecContext, func() (*slog.Attr, error) {
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
	logger   *stepLogger
	options  *stmtOptions
}

var _ driver.StmtQueryContext = (*stmtQueryContextWrapperImpl)(nil)

// QueryContext implements driver.StmtQueryContext.
func (s *stmtQueryContextWrapperImpl) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var rows driver.Rows
	err := ignoreAttr(lg.Step(ctx, &s.options.QueryContext, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.QueryContext(ctx, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return wrapRows(rows, s.logger, s.options.Rows), nil
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

// With driver.NamedValueChecker

type stmtNvcWrapper struct {
	stmtWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Stmt              = (*stmtNvcWrapper)(nil)
	_ driver.NamedValueChecker = (*stmtNvcWrapper)(nil)
)

type stmtExecContextNvcWrapper struct {
	stmtExecContextWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Stmt              = (*stmtExecContextNvcWrapper)(nil)
	_ driver.StmtExecContext   = (*stmtExecContextNvcWrapper)(nil)
	_ driver.NamedValueChecker = (*stmtExecContextNvcWrapper)(nil)
)

type stmtQueryContextNvcWrapper struct {
	stmtQueryContextWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Stmt              = (*stmtQueryContextNvcWrapper)(nil)
	_ driver.StmtQueryContext  = (*stmtQueryContextNvcWrapper)(nil)
	_ driver.NamedValueChecker = (*stmtQueryContextNvcWrapper)(nil)
)

type stmtContextNvcWrapper struct {
	stmtContextWrapper
	driver.NamedValueChecker
}

var (
	_ driver.Stmt              = (*stmtContextNvcWrapper)(nil)
	_ driver.StmtExecContext   = (*stmtContextNvcWrapper)(nil)
	_ driver.StmtQueryContext  = (*stmtContextNvcWrapper)(nil)
	_ driver.NamedValueChecker = (*stmtContextNvcWrapper)(nil)
)
