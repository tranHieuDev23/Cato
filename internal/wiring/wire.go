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

func InitializeLocalCato(
	filePath configs.ConfigFilePath,
	args utils.Arguments,
) (*app.LocalCato, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}

func InitializeDistributedHostCato(
	filePath configs.ConfigFilePath,
	args utils.Arguments,
) (*app.DistributedHostCato, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}

func InitializeDistributedWorkerCato(
	filePath configs.ConfigFilePath,
	args utils.Arguments,
) (*app.DistributedWorkerCato, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}
