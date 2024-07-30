package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg interface{})
	Error(msg interface{})
	Debug(msg interface{})
	Warning(msg interface{})
	Fatal(msg interface{})
}

type logger struct {
	logger *zap.Logger
}

func New(level string) (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(levelByString(level))

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &logger{
		logger: l,
	}, nil
}

func (l *logger) Info(msg interface{}) {
	l.logger.Info(toString(msg))
}

func (l *logger) Error(msg interface{}) {
	l.logger.Error(toString(msg))
}

func (l *logger) Debug(msg interface{}) {
	l.logger.Debug(toString(msg))
}

func (l *logger) Warning(msg interface{}) {
	l.logger.Warn(toString(msg))
}

func (l *logger) Fatal(msg interface{}) {
	l.logger.Fatal(toString(msg))
}

func levelByString(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return zapcore.FatalLevel
	case "error":
		return zapcore.ErrorLevel
	case "warning":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func toString(msg interface{}) string {
	return fmt.Sprintf("%v", msg)
}
