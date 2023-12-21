package main

import (
	"log"
	goHTTP "net/http"

	"github.com/tranHieuDev23/cato/internal/handlers/http"
	"github.com/tranHieuDev23/cato/web"
)

func main() {
	apiServerHandler := http.NewAPIServerHandler()
	spaHandler := http.NewSPAHandler(goHTTP.FS(web.StaticContent))
	server := http.NewServer(apiServerHandler, spaHandler)
	if err := server.Start(); err != nil {
		log.Println(err)
		panic(err)
	}
}
