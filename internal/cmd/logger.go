package cmd

import (
	"fmt"
	"io"
)

type stderrLogger struct {
	out io.Writer
}

func newStderrLogger(out io.Writer) *stderrLogger {
	return &stderrLogger{out: out}
}

func (l *stderrLogger) Debug(msg string, fields ...any) {
	l.write("DEBUG", msg, fields...)
}

func (l *stderrLogger) Info(msg string, fields ...any) {
	l.write("INFO", msg, fields...)
}

func (l *stderrLogger) Warn(msg string, fields ...any) {
	l.write("WARN", msg, fields...)
}

func (l *stderrLogger) Error(msg string, fields ...any) {
	l.write("ERROR", msg, fields...)
}

func (l *stderrLogger) write(level, msg string, fields ...any) {
	if l == nil || l.out == nil {
		return
	}
	if len(fields) == 0 {
		fmt.Fprintf(l.out, "[%s] %s\n", level, msg) //nolint:errcheck // Best-effort output
		return
	}
	fmt.Fprintf(l.out, "[%s] %s %v\n", level, msg, fields) //nolint:errcheck // Best-effort output
}
