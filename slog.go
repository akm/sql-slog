package sqlslog

import (
	"io"
	"log/slog"
	"os"
)

type slogOptions struct {
	slog.HandlerOptions
	handler     slog.Handler
	handlerFunc func(io.Writer, *slog.HandlerOptions) slog.Handler
	logWriter   io.Writer
}

func defaultSlogOptions() *slogOptions {
	return &slogOptions{
		handlerFunc:    NewTextHandler,
		logWriter:      os.Stdout,
		HandlerOptions: slog.HandlerOptions{},
	}
}

// Handler sets the slog.Handler to be used.
// If not set, the default is created by HandlerFunc, Writer, SlogOptions.
func Handler(handler slog.Handler) Option {
	return func(o *options) { o.SlogOptions.handler = handler }
}

// HandlerFunc sets the function to create the slog.Handler.
// If not set, the default is [NewTextHandler].
func HandlerFunc(handlerFunc func(io.Writer, *slog.HandlerOptions) slog.Handler) Option {
	return func(o *options) { o.SlogOptions.handlerFunc = handlerFunc }
}

// LogWriter sets the writer to be used for the slog.Handler.
// If not set, the default is os.Stdout.
func LogWriter(w io.Writer) Option {
	return func(o *options) { o.SlogOptions.logWriter = w }
}

// HandlerOptions sets the options to be used for the slog.Handler.
// If not set, the default is an empty [slog.HandlerOptions].
func HandlerOptions(opts *slog.HandlerOptions) Option {
	return func(o *options) {
		if opts == nil {
			opts = &slog.HandlerOptions{}
		}
		o.SlogOptions.HandlerOptions = *opts
	}
}

// AddSource sets whether to add the source to the log.
func AddSource(v bool) Option {
	return func(o *options) { o.SlogOptions.AddSource = v }
}

// LogLevel sets the log level to be used.
func LogLevel(v slog.Leveler) Option {
	return func(o *options) { o.SlogOptions.Level = v }
}

// ReplaceAttr sets the function to replace the attributes.
func LogReplaceAttr(f func([]string, slog.Attr) slog.Attr) Option {
	return func(o *options) { o.SlogOptions.ReplaceAttr = f }
}

// NewJSONHandler returns a new JSON handler using [slog.NewJSONHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return slog.NewJSONHandler(w, opts)
}

// NewTextHandler returns a new Text handler using [slog.NewTextHandler]
// with custom options for sqlslog.
// See [WrapHandlerOptions] for details on the options.
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	return slog.NewTextHandler(w, opts)
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
