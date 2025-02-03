// Package log implements custom log for logger.ContextLogger.
package log

import (
	"context"
	"io"

	F "github.com/sagernet/sing/common/format"
	"github.com/sagernet/sing/common/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ logger.ContextLogger = (*Logger)(nil)

type Logger struct {
	*zap.Logger
}

func New(ctx context.Context, writer io.Writer, level zapcore.Level, disableColor bool) *Logger {
	writeSyncer := zapcore.AddSync(writer)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	if disableColor {
		encoderConfig.EncodeLevel = defaultLevelEncoder
		encoderConfig.EncodeTime = defaultTimeEncoder
	} else {
		encoderConfig.EncodeLevel = colorfulLevelEncoder
		encoderConfig.EncodeTime = colorfulTimeEncoder
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	zapLogger := zap.New(core)

	return &Logger{
		Logger: zapLogger,
	}
}

func (l *Logger) Upstream() any {
	return l.Logger
}

func (l *Logger) Trace(args ...any) {
	l.TraceContext(context.Background(), args...)
}

func (l *Logger) Debug(args ...any) {
	l.DebugContext(context.Background(), args...)
}

func (l *Logger) Info(args ...any) {
	l.InfoContext(context.Background(), args...)
}

func (l *Logger) Warn(args ...any) {
	l.WarnContext(context.Background(), args...)
}

func (l *Logger) Error(args ...any) {
	l.ErrorContext(context.Background(), args...)
}

func (l *Logger) Fatal(args ...any) {
	l.FatalContext(context.Background(), args...)
}

func (l *Logger) Panic(args ...any) {
	l.PanicContext(context.Background(), args...)
}

func (l *Logger) TraceContext(_ context.Context, args ...any) {
	l.Logger.Debug(F.ToString(args...))
}

func (l *Logger) DebugContext(_ context.Context, args ...any) {
	l.Logger.Debug(F.ToString(args...))
}

func (l *Logger) InfoContext(_ context.Context, args ...any) {
	l.Logger.Info(F.ToString(args...))
}

func (l *Logger) WarnContext(_ context.Context, args ...any) {
	l.Logger.Warn(F.ToString(args...))
}

func (l *Logger) ErrorContext(_ context.Context, args ...any) {
	l.Logger.Error(F.ToString(args...))
}

func (l *Logger) FatalContext(_ context.Context, args ...any) {
	l.Logger.Fatal(F.ToString(args...))
}

func (l *Logger) PanicContext(_ context.Context, args ...any) {
	l.Logger.Panic(F.ToString(args...))
}
