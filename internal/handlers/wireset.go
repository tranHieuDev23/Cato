package handlers

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/handlers/http"
)

var WireSet = wire.NewSet(
	http.WireSet,
)
