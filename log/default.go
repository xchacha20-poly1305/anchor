package log

import (
	"context"
	"os"

	"go.uber.org/zap/zapcore"
)

var DefaultLog = New(context.Background(), os.Stderr, zapcore.InfoLevel, DisableColorFromEnv)

func Fatal(args ...any) {
	DefaultLog.Fatal(args...)
}

func Panic(args ...any) {
	DefaultLog.Fatal(args...)
}

func FatalContext(ctx context.Context, args ...any) {
	DefaultLog.FatalContext(ctx, args...)
}

func PanicContext(ctx context.Context, args ...any) {
	DefaultLog.FatalContext(ctx, args...)
}
