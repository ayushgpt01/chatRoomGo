package main

import (
	"fmt"      // formatting and printing values to the console.
	"log"      // logging messages to the console.
	"net/http" // Used for build HTTP servers and clients.
)

// Port we listen on.
const portNum string = ":8080"

func main() {
	fmt.Println()
	log.Println("Starting Web Sockets server")
	fmt.Println()

	// Create a Manager instance used to handle WebSocket Connections
	// manager := NewManager()

	// Registering our handler functions, and creating paths.
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// http.HandleFunc("/ws", manager.serveWS)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")
	fmt.Println()

	// Spinning up the server.
	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
