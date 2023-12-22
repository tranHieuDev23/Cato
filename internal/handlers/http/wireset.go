package http

import (
	"github.com/google/wire"
	"gitlab.com/pjrpc/pjrpc/v2"

	"github.com/tranHieuDev23/cato/internal/handlers/http/middlewares"
)

func InitializeMiddlewareList(
	authMiddleware middlewares.Auth,
) []pjrpc.Middleware {
	return []pjrpc.Middleware{
		pjrpc.Middleware(authMiddleware),
	}
}

var WireSet = wire.NewSet(
	middlewares.WireSet,
	InitializeMiddlewareList,
	NewAPIServerHandler,
	NewSPAHandler,
	NewServer,
)
