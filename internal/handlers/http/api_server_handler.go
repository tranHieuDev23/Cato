package http

import (
	"context"

	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc"
	"github.com/tranHieuDev23/cato/internal/handlers/http/rpc/rpcserver"
)

type apiServerHandler struct {
}

func NewAPIServerHandler() rpcserver.APIServer {
	return &apiServerHandler{}
}

func (apiServerHandler) Echo(ctx context.Context, in *rpc.EchoRequest) (*rpc.EchoResponse, error) {
	return &rpc.EchoResponse{
		Message: in.Message,
	}, nil
}
