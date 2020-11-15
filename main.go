package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Something went wrong 1,", err)
		return
	}

	fmt.Println("New client connected")
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("New message from client: %s", p)
	}
}

func main() {
	fmt.Println("Working...")

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":3000", nil)
	fmt.Println("Server is running")
}
