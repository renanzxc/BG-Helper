package log

import (
	"go.uber.org/zap"
)

func SetupLog() {
	cfgZap := zap.NewDevelopmentConfig()
	cfgZap.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfgZap.Encoding = "json"
	cfgZap.EncoderConfig.TimeKey = "date"
	cfgZap.EncoderConfig.LevelKey = "level"
	cfgZap.EncoderConfig.CallerKey = "caller"
	cfgZap.EncoderConfig.MessageKey = "message"

	logger, err := cfgZap.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
