package messaging

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/websocket"
)

// This factory is used in testing only.

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// CreateConnection - used to create a mock websocket connetion
func CreateConnection() (c websocket.Conn) {
	var s = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			c = *conn
			if err != nil {
				fmt.Printf("Something went wrong: %v", err.Error())
			}

			err = c.WriteJSON("test")
			if err != nil {
				fmt.Printf("Something went wrong: %v", err.Error())
			}
		}))
	defer s.Close()

	return
}
