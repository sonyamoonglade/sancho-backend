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
	builder.DisableStacktrace = !cfg.EnableStacktrace
	builder.Development = false
	builder.Encoding = "json"

	builder.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
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

func Get() *zap.Logger {
	return globalLogger
}
