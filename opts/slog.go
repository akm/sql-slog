package opts

import (
	"io"
	"log/slog"
)

// Logger sets the slog.Logger to be used.
// If not set, the default is slog.Default().
func Logger(logger *slog.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

// NewJSONHandler returns a new JSON handler using [slog.NewJSONHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	return slog.NewJSONHandler(w, WrapHandlerOptions(opts))
}

// NewTextHandler returns a new Text handler using [slog.NewTextHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) *slog.TextHandler {
	return slog.NewTextHandler(w, WrapHandlerOptions(opts))
}

// WrapHandlerOptions wraps the options with custom options for sqlslog.
// It merges ReplaceAttr functions with [ReplaceLevelAttr].
func WrapHandlerOptions(opts *slog.HandlerOptions) *slog.HandlerOptions {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	opts.ReplaceAttr = MergeReplaceAttrs(opts.ReplaceAttr, ReplaceLevelAttr)
	return opts
}

// ReplaceLevelAttr is a type of ReplaceAttr for [slog.HandlerOptions].
type ReplaceAttrFunc = func([]string, slog.Attr) slog.Attr

// MergeReplaceAttrs merges multiple [ReplaceAttrFunc] functions.
// If functions are nil or empty, it returns nil.
// If there is only one function, it returns that function.
// If there are multiple functions, it returns a merged function.
func MergeReplaceAttrs(funcs ...ReplaceAttrFunc) ReplaceAttrFunc {
	var valids []ReplaceAttrFunc
	for _, f := range funcs {
		if f != nil {
			valids = append(valids, f)
		}
	}
	if len(valids) == 0 {
		return nil
	}
	if len(valids) == 1 {
		return valids[0]
	}
	return func(group []string, a slog.Attr) slog.Attr {
		for _, f := range funcs {
			a = f(group, a)
		}
		return a
	}
}

// ReplaceLevelAttr replaces the log level as sqlslog.Level with its string representation.
func ReplaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	// https://go.dev/src/log/slog/example_custom_levels_test.go
	if a.Key == slog.LevelKey {
		level := Level(a.Value.Any().(slog.Level)) //nolint:forcetypeassert
		a.Value = slog.StringValue(level.String())
	}
	return a
}
