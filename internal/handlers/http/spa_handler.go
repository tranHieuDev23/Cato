package http

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tranHieuDev23/cato/web"
)

type SPAHandler http.Handler

type spaHandler struct {
	fileSystem http.FileSystem
	fileServer http.Handler
}

func NewSPAHandler() SPAHandler {
	return &spaHandler{
		fileSystem: http.FS(web.StaticContent),
		fileServer: http.FileServer(http.FS(web.StaticContent)),
	}
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	if strings.Contains(fileName, ".") {
		r.URL.Path = "dist/web/browser/" + r.URL.Path
		r.RequestURI = "/dist/web/browser" + r.RequestURI
		h.fileServer.ServeHTTP(w, r)
		return
	}

	indexFile, err := h.fileSystem.Open("/dist/web/browser/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "could not open embedded file", http.StatusInternalServerError)
		return
	}

	defer indexFile.Close()

	if _, err = io.Copy(w, indexFile); err != nil {
		http.Error(w, "could not copy file to response", http.StatusInternalServerError)
	}
}
