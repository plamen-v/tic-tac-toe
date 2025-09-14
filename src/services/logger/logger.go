package logger

import (
	"fmt"
	"time"

	"github.com/plamen-v/tic-tac-toe/src/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerService interface {
	Info(string, ...Field)
	Debug(string, ...Field)
	Error(string, ...Field)
	Sync() error
}

func NewLoggerService(appMode config.AppMode, logLevel config.LogLevel) (LoggerService, error) {
	logger := new(loggerService)
	if err := logger.initialize(appMode, logLevel); err != nil {
		return nil, err
	}
	return logger, nil
}

type Field struct {
	zapField zap.Field
}

func String(key, value string) Field {
	return Field{zapField: zap.String(key, value)}
}

func Int(key string, value int) Field {
	return Field{zapField: zap.Int(key, value)}
}

func Duration(key string, value time.Duration) Field {
	return Field{zapField: zap.Duration(key, value)}
}

type loggerService struct {
	worker *zap.Logger
}

func (l *loggerService) Info(msg string, fields ...Field) {
	l.worker.Info(msg, unwrap(fields)...)
}

func (l *loggerService) Debug(msg string, fields ...Field) {
	l.worker.Info(msg, unwrap(fields)...)
}

func (l *loggerService) Error(msg string, fields ...Field) {
	l.worker.Error(msg, unwrap(fields)...)
}

func (l *loggerService) Sync() error {
	return l.worker.Sync()
}

func (l *loggerService) initialize(appMode config.AppMode, logLevel config.LogLevel) error {
	var (
		level zapcore.Level
		cfg   zap.Config
		err   error
	)
	switch appMode {
	case config.DevelopmentAppMode:
		cfg = zap.NewDevelopmentConfig()
	case config.ProductionAppMode:
		cfg = zap.NewProductionConfig()
	default:
		return fmt.Errorf("unknown application mode '%s'", appMode)
	}

	if err = level.UnmarshalText([]byte(logLevel)); err != nil {
		return err
	}

	cfg.EncoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(level)

	l.worker, err = cfg.Build()
	return err
}

func unwrap(fields []Field) []zap.Field {
	zFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zFields[i] = f.zapField
	}
	return zFields
}
