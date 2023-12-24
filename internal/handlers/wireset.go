package handlers

import (
	"github.com/google/wire"

	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/internal/handlers/jobs"
)

var WireSet = wire.NewSet(
	http.WireSet,
	jobs.WireSet,
)
