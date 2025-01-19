package sqlslog

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestNewTextHandler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	h := NewTextHandler(buf, nil)
	if h == nil {
		t.Errorf("want not nil, got nil")
	}
}

func TestWrapHandlerOptions(t *testing.T) {
	t.Run("nil parameter", func(t *testing.T) {
		h := WrapHandlerOptions(nil)
		if h == nil {
			t.Errorf("want not nil, got nil")
		}
	})
	t.Run("one parameter", func(t *testing.T) {
		h := WrapHandlerOptions(&slog.HandlerOptions{})
		if h == nil {
			t.Errorf("want not nil, got nil")
		} else if h.ReplaceAttr == nil {
			t.Errorf("want not nil, got nil")
		}
	})
}

func TestMergeReplaceAttrs(t *testing.T) {
	t.Run("no parameter", func(t *testing.T) {
		f := MergeReplaceAttrs()
		if f != nil {
			t.Error("want nil, got not nil function")
		}
	})
	t.Run("nil parameter", func(t *testing.T) {
		f := MergeReplaceAttrs(nil, nil)
		if f != nil {
			t.Error("want nil, got not nil function")
		}
	})
	t.Run("one parameter", func(t *testing.T) {
		called := false
		f := MergeReplaceAttrs(func(group []string, a slog.Attr) slog.Attr {
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
		f := MergeReplaceAttrs(
			func(group []string, a slog.Attr) slog.Attr {
				return slog.String(a.Key, a.Value.String()+"+1")
			},
			func(group []string, a slog.Attr) slog.Attr {
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
