package sqlslog

import (
	"io"
	"log/slog"

	"github.com/akm/sql-slog/internal/opts"
)

// Logger sets the slog.Logger to be used.
// If not set, the default is slog.Default().
func Logger(logger *slog.Logger) Option {
	return opts.Logger(logger)
}

// NewJSONHandler returns a new JSON handler using [slog.NewJSONHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewJSONHandler(w io.Writer, options *slog.HandlerOptions) *slog.JSONHandler {
	return opts.NewJSONHandler(w, options)
}

// NewTextHandler returns a new Text handler using [slog.NewTextHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewTextHandler(w io.Writer, options *slog.HandlerOptions) *slog.TextHandler {
	return opts.NewTextHandler(w, options)
}
