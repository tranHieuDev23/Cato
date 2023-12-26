package utils

import (
	"os"
	"os/signal"
)

func BlockUntilSignal(signals ...os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, signals...)
	<-done
}
