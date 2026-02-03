package router

import (
	"net/http"
)

func handleViews(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("../../static")))
}
