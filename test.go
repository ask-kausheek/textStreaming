package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// struct of websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// incoming request from a any domain is allowed to connect
	CheckOrigin: func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}
func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		// print out that message for clarity
		fmt.Println("Received message:", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message:", err)
			return
		}

	}
}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
	}
	defer ws.Close()
	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("WebSocket Server started on :8080")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
