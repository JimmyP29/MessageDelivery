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
	//want := ""
	var s = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
			c = *conn
			if err != nil {
				//t.Fatal(err)
				fmt.Printf("Something went wrong: %v", err.Error())
			}

			// if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
			// 	//t.Errorf("resp.StatusCode = %q, want %q", got, want)
			// 	fmt.Printf("resp.StatusCode = %q, want %q", got, want)
			// }

			err = c.WriteJSON("test")
			if err != nil {
				//t.Fatal(err)
				fmt.Printf("Something went wrong: %v", err.Error())
			}
		}))
	defer s.Close()

	// var d = websocket.Dialer{}
	// c, resp, err := d.Dial("ws://"+s.Listener.Addr().String()+"/ws", nil)

	// if err != nil {
	// 	//t.Fatal(err)
	// 	fmt.Printf("Something went wrong: %v", err.Error())
	// }

	// if got, want := resp.StatusCode, http.StatusSwitchingProtocols; got != want {
	// 	//t.Errorf("resp.StatusCode = %q, want %q", got, want)
	// 	fmt.Printf("resp.StatusCode = %q, want %q", got, want)
	// }

	// err = c.WriteJSON("test")
	// if err != nil {
	// 	//t.Fatal(err)
	// 	fmt.Printf("Something went wrong: %v", err.Error())
	// }

	return
}

// func CreateClient() Client {
// 	return
// }
