package discovery

import (
	"context"
	"log"
	"log/slog"
)

// noopHandler implements Handler but does nothing.
type noopHandler struct{}

// Enabled always returns false so caller skips record construction.
func (noopHandler) Enabled(context.Context, slog.Level) bool { return false }

// Handle is never called (because Enabled is false) but must satisfy interface.
func (noopHandler) Handle(context.Context, slog.Record) error { return nil }

func (noopHandler) WithAttrs([]slog.Attr) slog.Handler { return noopHandler{} }
func (noopHandler) WithGroup(string) slog.Handler      { return noopHandler{} }

func (noopHandler) Write(p []byte) (int, error) { return len(p), nil }

// NewNoopLogger returns a Logger that drops all logs.
func NewNoopLogger() *slog.Logger {
	return slog.New(noopHandler{})
}

// NewNoopLogLogger returns a Logger that drops all logs.
func NewNoopLogLogger() *log.Logger {
	lg := log.New(noopHandler{}, "", 0)
	return lg
}

var NoopLogLogger = NewNoopLogLogger()
