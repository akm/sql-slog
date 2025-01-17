package sqlslog

import (
	"io"
	"log/slog"
)

func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	return slog.NewJSONHandler(w, WrapHandlerOptions(opts))
}

func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) *slog.JSONHandler {
	return slog.NewJSONHandler(w, WrapHandlerOptions(opts))
}

func WrapHandlerOptions(opts *slog.HandlerOptions) *slog.HandlerOptions {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	opts.ReplaceAttr = MergeReplaceAttrs(opts.ReplaceAttr, ReplaceLevelAttr)
	return opts
}

type ReplaceAttrFunc = func([]string, slog.Attr) slog.Attr

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
