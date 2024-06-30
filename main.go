package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

const portNum = ":8080"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendMessageEvent(evt Event) (Event, error) {
	payload, ok := evt.PAYLOAD.(map[string]interface{})
	if !ok {
		return evt, errors.New("invalid payload type")
	}

	msg := RecievedMessage{
		Message:    payload["message"].(string),
		SenderType: payload["senderType"].(string),
	}

	newMessage := Message{
		ID:         "1",
		Message:    msg.Message,
		SenderType: msg.SenderType,
		Username:   "Username",
		Status:     "Delivered",
	}
	log.Printf("Received message: %s", newMessage.Message)

	evt.PAYLOAD = newMessage
	evt.TYPE = "new_message"
	return evt, nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	for {
		var evt Event
		err := ws.ReadJSON(&evt)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		if evt.TYPE == "send_message" {
			newEvt, err := sendMessageEvent(evt)

			if err != nil {
				log.Printf("error: %v", err.Error())
				continue
			}

			err = ws.WriteJSON(newEvt)

			if err != nil {
				log.Printf("error: %v", err)
				break
			}
		}
	}
}

type RecievedMessage struct {
	Message    string `json:"message"`
	SenderType string `json:"senderType"`
}

type Message struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	SenderType string `json:"senderType"`
	Username   string `json:"username"`
	Status     string `json:"status"`
}

type Event struct {
	TYPE    string `json:"type"`
	PAYLOAD any    `json:"payload"`
}

func main() {
	fmt.Println()
	log.Println("Starting WebSockets server")
	fmt.Println()

	// Verify the current working directory and log it
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current working directory: %s", wd)

	// Log the absolute path of the static directory
	staticDir, err := filepath.Abs("./static")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Serving files from directory: %s", staticDir)

	// Verify that the static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Fatalf("Static directory does not exist: %s", staticDir)
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", handleConnections)

	log.Println("Started on port", portNum)
	fmt.Println("To close connection CTRL+C :-)")
	fmt.Println()

	err = http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
