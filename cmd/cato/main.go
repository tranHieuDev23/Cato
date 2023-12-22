package main

import (
	"log"

	"github.com/tranHieuDev23/cato/internal/wiring"
)

func main() {
	cato, cleanup, err := wiring.InitializeCato("")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer cleanup()
	if err := cato.Start(); err != nil {
		log.Println(err)
		panic(err)
	}
}
