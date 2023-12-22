package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Hash"),
	wire.FieldsOf(new(Config), "Token"),
)
