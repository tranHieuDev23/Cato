package configs

import (
	_ "embed"
)

//go:embed local.yaml
var DefaultConfigBytes []byte
