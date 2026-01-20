package logger

import "context"

type Logger interface {
	WithFields(fields map[string]interface{}) Logger
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

type ContextKey string

const ContextLoggerKey ContextKey = "logger"

func FromContext(ctx context.Context) Logger {
	val := ctx.Value(ContextLoggerKey)
	log, _ := val.(Logger)
	return log
}

func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}
