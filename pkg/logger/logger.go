package logger

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerFields map[string]interface{}

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
	WithFields(fields LoggerFields) Logger
}

type ZapLogger struct {
	logger *zap.Logger
	fields []zap.Field
}

func NewZapLogger(level string) (*ZapLogger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %v", err)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg, l.fields...)
}

func (l *ZapLogger) Warn(msg string) {
	l.logger.Warn(msg, l.fields...)
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg, l.fields...)
}

func (l *ZapLogger) Debug(msg string) {
	l.logger.Debug(msg, l.fields...)
}

func (l *ZapLogger) WithFields(fields LoggerFields) Logger {
	newLogger := *l
	for key, value := range fields {
		if subMap, ok := value.(LoggerFields); ok {
			newLogger.fields = append(newLogger.fields, zap.Any(key, subMap))
		} else {
			newLogger.fields = append(newLogger.fields, zap.Any(key, value))
		}
	}
	return &newLogger
}

func (l *ZapLogger) Close() {
	if err := l.logger.Sync(); err != nil {
		log.Println("logger sync error:", err)
	}
}
