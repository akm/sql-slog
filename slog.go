package sqlslog

import (
	"io"
	"log/slog"
)

// NewJSONHandler returns a new JSON handler by using [slog.NewJSONHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for the details of the options .
func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	return slog.NewJSONHandler(w, WrapHandlerOptions(opts))
}

// NewTextHandler returns a new Text handler by using [slog.NewTextHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for the details of the options .
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) *slog.TextHandler {
	return slog.NewTextHandler(w, WrapHandlerOptions(opts))
}

// WrapHandlerOptions wraps the options with custom options for sqlslog.
// It merge ReplaceAttr functions with [ReplaceLevelAttr].
func WrapHandlerOptions(opts *slog.HandlerOptions) *slog.HandlerOptions {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	opts.ReplaceAttr = MergeReplaceAttrs(opts.ReplaceAttr, ReplaceLevelAttr)
	return opts
}

// ReplaceLevelAttr is a type of ReplaceAttr of [slog.HandlerOptions].
type ReplaceAttrFunc = func([]string, slog.Attr) slog.Attr

// ReplaceLevelAttr merges the [ReplaceAttrFunc] functions.
// If functions are nil or empty, it returns nil.
// If functions are only one, it returns the function.
// If functions are more than one, it returns the merged function.
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
