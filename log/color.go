package log

import (
	"os"
	"time"

	"go.uber.org/zap/zapcore"
)

// DisableColorFromEnv checks environment to decides whether we should use color.
//
// https://no-color.org/
var DisableColorFromEnv = os.Getenv("NO_COLOR") != ""

// ANSI color codes as constants
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
)

func colorfulLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = ColorCyan
	case zapcore.InfoLevel:
		color = ColorGreen
	case zapcore.WarnLevel:
		color = ColorYellow
	case zapcore.ErrorLevel:
		color = ColorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = ColorMagenta
	default:
		color = ColorWhite
	}
	enc.AppendString(color + level.CapitalString() + ColorReset)
}

func colorfulTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(ColorYellow + t.Format(time.RFC3339) + ColorReset)
}

func defaultLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.CapitalString())
}

func defaultTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339))
}
