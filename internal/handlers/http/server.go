package http

import (
	"log"
	"net/http"

	"gitlab.com/pjrpc/pjrpc/v2"

	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
)

type Server interface {
	Start() error
}

type server struct {
	apiServerHandler rpcserver.APIServer
	middlewareList   []pjrpc.Middleware
	spaHandler       http.Handler
}

func NewServer(
	apiServerHandler rpcserver.APIServer,
	middlewareList []pjrpc.Middleware,
	spaHandler SPAHandler,
) Server {
	return &server{
		apiServerHandler: apiServerHandler,
		spaHandler:       spaHandler,
	}
}

func (s server) Start() error {
	srv := pjrpc.NewServerHTTP()
	srv.SetLogger(log.Writer()) // Server can write body close errors and panics in handlers.

	rpcserver.RegisterAPIServer(srv, s.apiServerHandler, s.middlewareList...)

	mux := http.NewServeMux()

	// Be careful about last slash. Requests will be sent with last slash by swagger UI.
	// Recomend just to strip last slash in the router.
	mux.Handle("/api/", srv)

	// Serving static files
	mux.Handle("/", s.spaHandler)

	log.Println("Starting rpc server on :8080")

	return http.ListenAndServe(":8080", mux)

}
