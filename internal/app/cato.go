package app

import (
	"go.uber.org/zap"

	"github.com/tranHieuDev23/cato/internal/handlers/http"
)

type Cato interface {
	Start() error
}

type cato struct {
	httpServer http.Server
	logger     *zap.Logger
}

func NewCato(
	httpServer http.Server,
	logger *zap.Logger,
) Cato {
	return &cato{}
}

func (c cato) Start() error {
	return c.httpServer.Start()
}
