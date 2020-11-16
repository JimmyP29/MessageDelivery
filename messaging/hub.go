package messaging

import (
	"log"
	"net"
)

type hub struct {
	clients  map[uint64]*client
	messages chan message
}

// NewHub Returns a hub reference
func NewHub() *hub {
	return &hub{
		clients:  make(map[uint64]*client),
		messages: make(chan message),
	}
}

func (h *hub) NewClient(conn net.Conn) {
	log.Printf("New Client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn: conn,
		//userID: uint64,// random number?
		messages: h.messages,
	}
}
