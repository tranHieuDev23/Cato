package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/pjrpc/pjrpc/v2"
	"go.uber.org/zap"

	"github.com/pkg/browser"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
	"github.com/tranHieuDev23/cato/internal/utils"
)

type Server interface {
	Start() error
}

type server struct {
	apiServerHandler    rpcserver.APIServer
	pjrpcMiddlewareList []pjrpc.Middleware
	httpMiddlewareList  []func(http.Handler) http.Handler
	spaHandler          http.Handler
	logger              *zap.Logger
	httpConfig          configs.HTTP
	appArguments        utils.Arguments
}

func NewServer(
	apiServerHandler rpcserver.APIServer,
	middlewareList []pjrpc.Middleware,
	httpMiddlewareList []func(http.Handler) http.Handler,
	spaHandler SPAHandler,
	logger *zap.Logger,
	httpConfig configs.HTTP,
	appArguments utils.Arguments,
) Server {
	return &server{
		apiServerHandler:    apiServerHandler,
		pjrpcMiddlewareList: middlewareList,
		httpMiddlewareList:  httpMiddlewareList,
		spaHandler:          spaHandler,
		logger:              logger,
		httpConfig:          httpConfig,
		appArguments:        appArguments,
	}
}

func (s server) Start() error {
	srv := pjrpc.NewServerHTTP()
	srv.OnPanic = func(ctx context.Context, err error) *pjrpc.ErrorResponse {
		utils.LoggerWithContext(ctx, s.logger).With(zap.Error(err)).Error("panic occurred")
		return srv.DefaultRestoreOnPanic(ctx, err)
	}

	rpcserver.RegisterAPIServer(srv, s.apiServerHandler, s.pjrpcMiddlewareList...)

	var apiHandler http.Handler = srv
	for i := range s.httpMiddlewareList {
		apiHandler = s.httpMiddlewareList[i](apiHandler)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/", s.spaHandler)

	httpServer := &http.Server{
		Addr:              s.httpConfig.Address,
		Handler:           mux,
		ReadHeaderTimeout: time.Minute,
	}

	if !s.appArguments.NoBrowser {
		time.AfterFunc(time.Second, func() {
			_ = browser.OpenURL(fmt.Sprintf("http://%s", s.httpConfig.Address))
		})
	}

	s.logger.With(zap.String("address", s.httpConfig.Address)).Info("starting http server")
	return httpServer.ListenAndServe()
}
