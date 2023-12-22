package middlewares

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAuth,
)
