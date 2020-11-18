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

var hub = messaging.NewHub(
	make([]messaging.Client, 100),
	make([]messaging.Subscription, 100),
)

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

	// Prior to a new client connecting, they are assigned an id
	c := messaging.NewClient(hub.AssignUserID(), conn)

	// The client is first added to the hub and then subscribed to the topic
	hub.AddClient(*c)
	hub.Subscribe(*c)

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
