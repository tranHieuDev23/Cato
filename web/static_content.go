package web

import "embed"

//go:embed dist/web/browser/*
var StaticContent embed.FS
