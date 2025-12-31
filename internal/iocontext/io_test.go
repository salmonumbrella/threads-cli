// internal/iocontext/io_test.go
package iocontext

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestDefaultIO(t *testing.T) {
	io := DefaultIO()
	if io.Out != os.Stdout {
		t.Error("expected stdout")
	}
	if io.ErrOut != os.Stderr {
		t.Error("expected stderr")
	}
	if io.In != os.Stdin {
		t.Error("expected stdin")
	}
}

func TestWithIO(t *testing.T) {
	var buf bytes.Buffer
	io := &IO{Out: &buf, ErrOut: &buf, In: strings.NewReader("test")}
	ctx := WithIO(context.Background(), io)

	got := GetIO(ctx)
	if got != io {
		t.Error("expected injected IO")
	}
}

func TestGetIO_FallsBackToDefault(t *testing.T) {
	ctx := context.Background()
	io := GetIO(ctx)
	if io.Out != os.Stdout {
		t.Error("expected fallback to stdout")
	}
}

func TestHasIO(t *testing.T) {
	ctx := context.Background()
	if HasIO(ctx) {
		t.Error("expected no IO in empty context")
	}

	ctx = WithIO(ctx, &IO{})
	if !HasIO(ctx) {
		t.Error("expected IO after injection")
	}
}
