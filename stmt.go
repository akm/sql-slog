package sqlslog

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log/slog"
)

type StmtOptions struct {
	Close        *StepOptions
	Exec         *StepOptions
	Query        *StepOptions
	ExecContext  *StepOptions
	QueryContext *StepOptions

	Rows *RowsOptions
}

func DefaultStmtOptions(formatter StepLogMsgFormatter) *StmtOptions {
	return &StmtOptions{
		Close:        DefaultStepOptions(formatter, "Stmt.Close", LevelInfo),
		Exec:         DefaultStepOptions(formatter, "Stmt.Exec", LevelInfo),
		Query:        DefaultStepOptions(formatter, "Stmt.Query", LevelInfo),
		ExecContext:  DefaultStepOptions(formatter, "Stmt.ExecContext", LevelInfo),
		QueryContext: DefaultStepOptions(formatter, "Stmt.QueryContext", LevelInfo),
		Rows:         DefaultRowsOptions(formatter),
	}
}

type stmtWrapper struct {
	original driver.Stmt
	logger   *StepLogger
	options  *StmtOptions
}

var _ driver.Stmt = (*stmtWrapper)(nil)

// Close implements driver.Stmt.
func (s *stmtWrapper) Close() error {
	return IgnoreAttr(s.logger.StepWithoutContext(s.options.Close, WithNilAttr(s.original.Close)))
}

// Exec implements driver.Stmt.
func (s *stmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := IgnoreAttr(lg.StepWithoutContext(s.options.Exec, func() (*slog.Attr, error) {
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
	err := IgnoreAttr(lg.StepWithoutContext(s.options.Query, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.Query(args) //nolint:staticcheck
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, s.logger, s.options.Rows), nil
}

type stmtExecContextWrapperImpl struct {
	original driver.StmtExecContext
	logger   *StepLogger
	options  *StmtOptions
}

var _ driver.StmtExecContext = (*stmtExecContextWrapperImpl)(nil)

// ExecContext implements driver.StmtExecContext.
func (s *stmtExecContextWrapperImpl) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var result driver.Result
	err := IgnoreAttr(lg.Step(ctx, s.options.ExecContext, func() (*slog.Attr, error) {
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
	logger   *StepLogger
	options  *StmtOptions
}

var _ driver.StmtQueryContext = (*stmtQueryContextWrapperImpl)(nil)

// QueryContext implements driver.StmtQueryContext.
func (s *stmtQueryContextWrapperImpl) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	lg := s.logger.With(slog.String("args", fmt.Sprintf("%+v", args)))
	var rows driver.Rows
	err := IgnoreAttr(lg.Step(ctx, s.options.QueryContext, func() (*slog.Attr, error) {
		var err error
		rows, err = s.original.QueryContext(ctx, args)
		return nil, err
	}))
	if err != nil {
		return nil, err
	}
	return WrapRows(rows, s.logger, s.options.Rows), nil
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

func WrapStmt(original driver.Stmt, logger *StepLogger, options *StmtOptions) driver.Stmt {
	if original == nil {
		return nil
	}
	stmtWrapper := stmtWrapper{original: original, logger: logger, options: options}

	stmtExec, withExecContext := original.(driver.StmtExecContext)
	stmtQuery, withQueryContext := original.(driver.StmtQueryContext)
	if withExecContext && withQueryContext {
		return &stmtContextWrapper{
			stmtWrapper:                 stmtWrapper,
			stmtExecContextWrapperImpl:  stmtExecContextWrapperImpl{original: stmtExec, logger: logger, options: options},
			stmtQueryContextWrapperImpl: stmtQueryContextWrapperImpl{original: stmtQuery, logger: logger, options: options},
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
