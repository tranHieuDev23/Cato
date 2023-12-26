package app

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHost,
	NewWorker,
)
