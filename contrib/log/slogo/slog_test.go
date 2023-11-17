package slogo

import (
	"context"
	"log/slog"
	"testing"
)

// testHandler is a Handler just for testing that calls optional hooks on each method.
type testHandler struct {
	fnWithAttrs func(attrs []slog.Attr)
	fnWithGroup func(name string)
}

var _ slog.Handler = &testHandler{}

func (h *testHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return true
}

func (h *testHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (h *testHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.fnWithAttrs != nil {
		h.fnWithAttrs(attrs)
	}
	out := *h
	return &out
}

func (h *testHandler) WithGroup(name string) slog.Handler {
	if h.fnWithGroup != nil {
		h.fnWithGroup(name)
	}
	out := *h
	return &out
}

func TestWithName(t *testing.T) {
	calledWithAttrs := 0
	nameInput := "name"

	handler := &testHandler{}
	handler.fnWithAttrs = func(attrs []slog.Attr) {
		calledWithAttrs++

		name := attrs[0].Value.String()
		if name != nameInput {
			t.Errorf("unexpected name input, got %q", name)
		}
	}
	log := slog.New(handler)

	out := WithName(log, nameInput)

	if calledWithAttrs != 1 {
		t.Errorf("expected handler.WithAttrs() to be called once, got %d", calledWithAttrs)
	}

	if p, _ := out.Handler().(*testHandler); p == handler {
		t.Errorf("expected output to be different from input, got in=%p, out=%p", handler, p)
	}
}
