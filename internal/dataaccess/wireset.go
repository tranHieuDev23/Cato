package dataaccess

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/dataaccess/cache"
	"github.com/tranHieuDev23/cato/internal/dataaccess/cato"
	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
)

var WireSet = wire.NewSet(
	cato.WireSet,
	db.WireSet,
	cache.WireSet,
)
