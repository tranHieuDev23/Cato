package dataaccess

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/dataaccess/db"
)

var WireSet = wire.NewSet(
	db.WireSet,
)
