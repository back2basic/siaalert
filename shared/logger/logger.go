package logger

import (
    "os"
    "sync"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    loggerInstance *zap.Logger
    once           sync.Once
)

// GetLogger initializes and returns a singleton logger instance
func GetLogger(file string) *zap.Logger {
    once.Do(func() {
        // Custom encoder configuration for colorized logging
        consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
            TimeKey:        "time",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "msg",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    colorizedLevelEncoder, // Add colors to log levels
            EncodeTime:     zapcore.ISO8601TimeEncoder,
            EncodeDuration: zapcore.StringDurationEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
        })

        // Output to both stdout and a log file
        consoleWriter := zapcore.Lock(os.Stdout)
        fileWriter, _ := os.Create(file)

        core := zapcore.NewTee(
            zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel), // Colorized stdout
            zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
                TimeKey:        "time",
                LevelKey:       "level",
                MessageKey:     "msg",
                StacktraceKey:  "stacktrace",
                EncodeTime:     zapcore.ISO8601TimeEncoder,
                EncodeLevel:    zapcore.LowercaseLevelEncoder,
            }), zapcore.AddSync(fileWriter), zapcore.DebugLevel), // JSON to log file
        )

        loggerInstance = zap.New(core)
    })
    return loggerInstance
}

// Sync flushes any buffered log entries
func Sync() {
    if loggerInstance != nil {
        _ = loggerInstance.Sync()
    }
}

// colorizedLevelEncoder applies colors to log levels
func colorizedLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
    var levelColor string
    switch level {
    case zapcore.DebugLevel:
        levelColor = "\033[34m" // Blue
		case zapcore.InfoLevel:
				levelColor = "\033[32m" // Green
		case zapcore.WarnLevel:
				levelColor = "\033[33m" // Yellow
		case zapcore.ErrorLevel:
				levelColor = "\033[31m" // Red
		default:
				levelColor = "\033[0m" // Reset
		}
		enc.AppendString(levelColor + level.CapitalString() + "\033[0m") // Reset
}