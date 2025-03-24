package logger

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func init() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.OutputPaths = []string{"logs.log"}
	loggerConfig.ErrorOutputPaths = []string{"logs.log"}
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	Log = logger.Sugar()
}
