package main

import (
	"MessageDelivery/messaging"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = messaging.NewHub(make([]messaging.Client, 100))

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("New client connected")

	c := messaging.NewClient(hub.AssignUserID(), conn)

	hub.AddClient(*c)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		hub.HandleReceiveMessage(*c, messageType, p)
		// if err := conn.WriteMessage(messageType, p); err != nil {
		// 	log.Println(err)
		// 	return
		// }

		// fmt.Printf("New message from client: %s \n", p)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "html")
	})
	http.HandleFunc("/ws", websocketHandler)
	http.ListenAndServe(":8888", nil)
	fmt.Println("Server running...")
}
