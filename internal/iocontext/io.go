// Package iocontext provides IO stream abstraction for testability.
// It allows injecting custom io.Reader/io.Writer into context,
// enabling tests to capture output and provide input without
// using os.Stdout/Stderr/Stdin directly.
package iocontext

import (
	"context"
	"io"
	"os"
)

// IO holds input/output streams for commands
type IO struct {
	Out    io.Writer // stdout
	ErrOut io.Writer // stderr
	In     io.Reader // stdin
}

type contextKey struct{}

// DefaultIO returns IO using os.Std streams
func DefaultIO() *IO {
	return &IO{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
		In:     os.Stdin,
	}
}

// WithIO injects IO into context
func WithIO(ctx context.Context, io *IO) context.Context {
	return context.WithValue(ctx, contextKey{}, io)
}

// GetIO retrieves IO from context, falling back to defaults
func GetIO(ctx context.Context) *IO {
	if io, ok := ctx.Value(contextKey{}).(*IO); ok && io != nil {
		return io
	}
	return DefaultIO()
}

// HasIO checks if IO is in context
func HasIO(ctx context.Context) bool {
	_, ok := ctx.Value(contextKey{}).(*IO)
	return ok
}
