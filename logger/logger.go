package logger

import (
	"go.uber.org/zap"
)

var globalLogger *zap.Logger

type Config struct {
	Out              []string
	Strict           bool
	Production       bool
	EnableStacktrace bool
}

func NewLogger(cfg Config) error {

	builder := zap.NewProductionConfig()
	builder.DisableStacktrace = true
	builder.Development = false
	builder.Encoding = "JSON"
	builder.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if cfg.Strict {
		builder.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	}
	if cfg.Production {
		builder.OutputPaths = cfg.Out
	}

	logger, err := builder.Build()
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

func GetLogger() *zap.Logger {
	return globalLogger
}
