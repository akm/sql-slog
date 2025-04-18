package sqlslog

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestNewTextHandler(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer(nil)
	h := NewTextHandler(buf, nil)
	if h == nil {
		t.Errorf("want not nil, got nil")
	}
}

func TestWrapHandlerOptions(t *testing.T) {
	t.Parallel()
	t.Run("nil parameter", func(t *testing.T) {
		t.Parallel()
		h := WrapHandlerOptions(nil)
		if h == nil {
			t.Errorf("want not nil, got nil")
		}
	})
	t.Run("one parameter", func(t *testing.T) {
		t.Parallel()
		h := WrapHandlerOptions(&slog.HandlerOptions{})
		if h == nil {
			t.Errorf("want not nil, got nil")
		} else if h.ReplaceAttr == nil {
			t.Errorf("want not nil, got nil")
		}
	})
}

func TestMergeReplaceAttrs(t *testing.T) {
	t.Parallel()
	t.Run("no parameter", func(t *testing.T) {
		t.Parallel()
		f := MergeReplaceAttrs()
		if f != nil {
			t.Error("want nil, got not nil function")
		}
	})
	t.Run("nil parameter", func(t *testing.T) {
		t.Parallel()
		f := MergeReplaceAttrs(nil, nil)
		if f != nil {
			t.Error("want nil, got not nil function")
		}
	})
	t.Run("one parameter", func(t *testing.T) {
		t.Parallel()
		called := false
		f := MergeReplaceAttrs(func(_ []string, a slog.Attr) slog.Attr {
			called = true
			return a
		})
		if f == nil {
			t.Errorf("want not nil, got nil")
		}
		r := f([]string{}, slog.String("key", "value"))
		if !called {
			t.Errorf("want called, got not called")
		}
		if r.Key != "key" {
			t.Errorf("want key as r.Key, got %s", r.Key)
		}
		if r.Value.Kind() != slog.KindString || r.Value.String() != "value" {
			t.Errorf("want value as r.Value, got %v", r.Value)
		}
	})
	t.Run("two parameters", func(t *testing.T) {
		t.Parallel()
		f := MergeReplaceAttrs(
			func(_ []string, a slog.Attr) slog.Attr {
				return slog.String(a.Key, a.Value.String()+"+1")
			},
			func(_ []string, a slog.Attr) slog.Attr {
				return slog.String(a.Key, a.Value.String()+"+2")
			},
		)
		if f == nil {
			t.Errorf("want not nil, got nil")
		}
		r := f([]string{}, slog.String("key", "value"))
		if r.Key != "key" {
			t.Errorf("want key as r.Key, got %s", r.Key)
		}
		if r.Value.Kind() != slog.KindString || r.Value.String() != "value+1+2" {
			t.Errorf("want value as r.Value, got %v", r.Value)
		}
	})
}

func TestHandlerOptions(t *testing.T) {
	t.Parallel()
	t.Run("nil parameter", func(t *testing.T) {
		t.Parallel()
		o := HandlerOptions(nil)
		if o == nil {
			t.Errorf("want not nil, got nil")
		}
		opts := &options{SlogOptions: defaultSlogOptions()}
		o(opts)
	})
	t.Run("one parameter", func(t *testing.T) {
		t.Parallel()
		o := HandlerOptions(&slog.HandlerOptions{})
		if o == nil {
			t.Errorf("want not nil, got nil")
		}
		opts := &options{SlogOptions: defaultSlogOptions()}
		o(opts)
	})
}

func TestAddSource(t *testing.T) {
	t.Parallel()
	o := AddSource(true)
	if o == nil {
		t.Errorf("want not nil, got nil")
	}
	opts := &options{SlogOptions: defaultSlogOptions()}
	o(opts)
}
