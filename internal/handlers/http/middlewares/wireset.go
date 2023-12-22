package middlewares

import (
	"net/http"

	"github.com/google/wire"
	"gitlab.com/pjrpc/pjrpc/v2"
)

func InitializePJRPCMiddlewareList() []pjrpc.Middleware {
	return []pjrpc.Middleware{}
}

func InitalizeHTTPMiddlewareList(
	httpAuthMiddleware HTTPAuth,
) []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		httpAuthMiddleware,
	}
}

var WireSet = wire.NewSet(
	InitializePJRPCMiddlewareList,
	InitalizeHTTPMiddlewareList,
	NewHTTPAuth,
)
