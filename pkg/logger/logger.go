package logger

import (
	"context"
	"encoding/json"
	"fmt"
)

type L struct{}

func New() L {
	return L{}
}

func (l L) Debug(ctx context.Context, format string, args ...any) {
	l.emit(ctx, "DEBUG", format, args...)
}

func (l L) Info(ctx context.Context, format string, args ...any) {
	l.emit(ctx, "INFO", format, args...)
}

func (l L) Warn(ctx context.Context, format string, args ...any) {
	l.emit(ctx, "WARN", format, args...)
}

func (l L) Error(ctx context.Context, format string, args ...any) {
	l.emit(ctx, "ERROR", format, args...)
}

func (l L) Critical(ctx context.Context, format string, args ...any) {
	l.emit(ctx, "CRITICAL", format, args...)
}

func (l L) emit(_ context.Context, severity, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	sl := struct {
		Message  string `json:"message"`
		Severity string `json:"severity"`
	}{
		Message:  msg,
		Severity: severity,
	}
	json, err := json.Marshal(sl)
	if err != nil {
		fmt.Printf("Failed to marshal log %s %q: %v\n", severity, msg, err)
		return
	}
	fmt.Printf("%s\n", json)
}
