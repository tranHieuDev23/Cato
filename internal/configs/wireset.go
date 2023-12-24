package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Database"),
	wire.FieldsOf(new(Config), "Auth"),
	wire.FieldsOf(new(Config), "HTTP"),
	wire.FieldsOf(new(Config), "Logic"),
	wire.FieldsOf(new(Auth), "Hash"),
	wire.FieldsOf(new(Auth), "Token"),
)
