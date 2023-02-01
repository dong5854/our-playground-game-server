package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	var err error

	logger, err = myDebugConfig().Build()
	if err != nil {
		panic(err)
	}
}

func myDebugConfig() zap.Config {
	encodeCfg := zapcore.EncoderConfig{
		StacktraceKey: "stacktrace",
		MessageKey:    "message",
		CallerKey:     "caller",
		EncodeCaller:  zapcore.ShortCallerEncoder,
		TimeKey:       "timestamp",
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		LevelKey:      "level",
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
	}
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encodeCfg,
	}
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Sync() {
	logger.Sync()
}
