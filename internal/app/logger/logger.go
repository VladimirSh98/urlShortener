package logger

import (
	"go.uber.org/zap"
)

func Initialize() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)
	sugar := *logger.Sugar()
	sugar.Infow("Logger initialized")
	return logger, nil
}
