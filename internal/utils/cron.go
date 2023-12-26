package utils

import (
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

func InitializeGoCronScheduler(logger *zap.Logger) (gocron.Scheduler, func(), error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, nil, err
	}

	cleanupFunc := func() {
		err = scheduler.Shutdown()
		if err != nil {
			logger.With(zap.Error(err)).Error("failed to shutdown scheduler")
		}
	}

	return scheduler, cleanupFunc, nil
}
