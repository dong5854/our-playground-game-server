package logger

import (
	"encoding/json"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var cfg zap.Config
	var err error

	JSONConf := []byte(`{
		"level": "debug",
		"encoding" : "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig" : {
			"messageKey" : "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)

	if err := json.Unmarshal(JSONConf, &cfg); err != nil {
		panic(err)
	}

	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Sync() {
	logger.Sync()
}
