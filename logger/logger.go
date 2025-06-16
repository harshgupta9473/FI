package logger

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

func NewLogger() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{
		logger: log,
	}, nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(err error, msg string, fields ...zap.Field) {
	l.logger.Error(msg, append(fields, zap.Error(err))...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Sync() {
	_ = l.logger.Sync()
}
