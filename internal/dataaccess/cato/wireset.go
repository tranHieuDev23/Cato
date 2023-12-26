package cato

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewHTTPClientWithAuthToken,
	InitializeAuthenticatedClient,
)
