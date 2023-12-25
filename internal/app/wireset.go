package app

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewLocalCato,
	NewDistributedHostCato,
	NewDistributedWorkerCato,
)
