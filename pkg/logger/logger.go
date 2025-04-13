package logger

import (
	"go.uber.org/zap"
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

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
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
		if subMap, ok := value.(LoggerFields); ok{
			newLogger.fields = append(newLogger.fields, zap.Any(key, subMap))
		} else {
			newLogger.fields = append(newLogger.fields, zap.Any(key, value))
		}
	}
	return &newLogger
}

func (l *ZapLogger) Close() {
	l.logger.Sync()
}

