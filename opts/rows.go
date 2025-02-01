package opts

import (
	"errors"
	"io"
	"log/slog"
)

type RowsOptions struct {
	Close         *StepOptions
	Next          *StepOptions
	NextResultSet *StepOptions
}

func DefaultRowsOptions(formatter StepLogMsgFormatter) *RowsOptions {
	return &RowsOptions{
		Close:         DefaultStepOptions(formatter, "Rows.Close", LevelDebug),
		Next:          DefaultStepOptions(formatter, "Rows.Next", LevelDebug, HandleRowsNextError),
		NextResultSet: DefaultStepOptions(formatter, "Rows.NextResultSet", LevelDebug),
	}
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
