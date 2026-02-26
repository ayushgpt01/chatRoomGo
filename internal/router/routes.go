package router

import (
	"log"
	"net/http"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/ws"
)

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Recieved Request from URL %s\n", r.URL.Path)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("Took %v time to respond to URL %s\n", t2.Sub(t1), r.URL.Path)
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered in f %v\n", r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func HandleRoutes(wsHandler *ws.Wshandler, authService *auth.AuthService, roomService *room.RoomService) http.Handler {
	log.Printf("Setting up routes...")

	mux := http.NewServeMux()

	handleAPIRoutes(mux, authService, roomService)
	handleViews(mux)
	mux.Handle("/ws", wsHandler)

	var handler http.Handler = mux
	handler = recoverMiddleware(handler)
	handler = metricsMiddleware(handler)

	return handler
}
