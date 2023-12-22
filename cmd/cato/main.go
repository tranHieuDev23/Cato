package main

import (
	"log"
	"os"

	"github.com/tranHieuDev23/cato/internal/configs"
	"github.com/tranHieuDev23/cato/internal/wiring"
)

func main() {
	configFilePath := ""
	if len(os.Args) == 2 {
		configFilePath = os.Args[1]
	}

	cato, cleanup, err := wiring.InitializeCato(configs.ConfigFilePath(configFilePath))
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
