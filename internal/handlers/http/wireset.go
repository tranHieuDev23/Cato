package http

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
)

var WireSet = wire.NewSet(
	middlewares.WireSet,
	NewAPIServerHandler,
	NewSPAHandler,
	NewServer,
)
