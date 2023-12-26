package utils

import (
	"github.com/robfig/cron/v3"
)

func InitializeCron() (*cron.Cron, func()) {
	c := cron.New()
	cleanupFunc := func() {
		<-c.Stop().Done()
	}

	return c, cleanupFunc
}
