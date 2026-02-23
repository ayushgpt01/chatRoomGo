package router

import (
	"net/http"
	"os"
	"path/filepath"
)

func handleViews(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("../../static", r.URL.Path)

		// If the file exists, serve it. If not, serve index.html (for the SPA)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, "../../static/index.html")
			return
		}
		http.FileServer(http.Dir("../../static")).ServeHTTP(w, r)
	})
}
