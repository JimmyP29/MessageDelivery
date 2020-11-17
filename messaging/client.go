package messaging

import "github.com/gorilla/websocket"

// Client - represents a user in the system
type Client struct {
	userID     uint64
	connection *websocket.Conn
}

// NewClient creates a new client
func NewClient(userID uint64, connection *websocket.Conn) *Client {
	return &Client{userID, connection}
}
