package router

import (
	"net/http"
	"os"
	"path/filepath"
)

func handleViews(mux *http.ServeMux) {
	staticDir := os.Getenv("STATIC_FILES_PATH")
	if staticDir == "" {
		staticDir = "./static"
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fullPath := filepath.Join(staticDir, r.URL.Path)

		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) || info.IsDir() {
			// Serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
	})
}
