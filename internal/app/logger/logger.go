package logger

import (
	"go.uber.org/zap"
)

func Initialize() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	sugar := *logger.Sugar()
	sugar.Infow("Logger initialized")
	return nil
}
