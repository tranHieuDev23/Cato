//go:build wireinject
// +build wireinject

//
//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/app"
	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/dataaccess"
	"github.com/tranHieuDev23/cato/internal/handlers"
	"github.com/tranHieuDev23/cato/internal/logic"
	"github.com/tranHieuDev23/cato/internal/utils"
)

var WireSet = wire.NewSet(
	utils.WireSet,
	configs.WireSet,
	dataaccess.WireSet,
	logic.WireSet,
	handlers.WireSet,
	app.WireSet,
)

func InitializeHost(
	filePath configs.ConfigFilePath,
	args utils.Arguments,
) (*app.Host, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}

func InitializeWorker(
	filePath configs.ConfigFilePath,
	args utils.Arguments,
) (*app.Worker, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}
