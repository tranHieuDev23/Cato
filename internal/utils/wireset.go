package utils

import "github.com/google/wire"

var WireSet = wire.NewSet(
	InitializeLogger,
)
