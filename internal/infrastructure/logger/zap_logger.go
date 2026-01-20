package logger

import (
	"os"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(cfg *config.Config) *zapLogger {
	level := zapcore.DebugLevel
	if cfg.IsProduction() {
		level = zapcore.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "message",
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// Configurando o output para console e nível de log
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level,
	)

	logger := zap.New(core)

	return &zapLogger{logger: logger}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) WithFields(fields map[string]any) logger.Logger {
	zapFields := []zapcore.Field{}
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	newLogger := l.logger.With(zapFields...)

	return &zapLogger{logger: newLogger}
}
