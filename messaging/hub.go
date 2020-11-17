package messaging

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Hub - used to control message flow to clients
type Hub struct {
	clients []Client
}

// NewHub creates a new hub
func NewHub(clients []Client) *Hub {
	return &Hub{
		clients,
	}
}

// AssignUserID - creates new uint64 userID using RNG
func (h *Hub) AssignUserID() uint64 {
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Intn(1000))
}

// AddClient - appends new client to client slice
func (h *Hub) AddClient(client Client) *Hub {
	h.clients = append(h.clients, client)

	fmt.Printf("Client %v has been added to the hub \n", client.userID)
	//fmt.Printf("hub: %+v \n", h.clients)
	return h
}

// HandleReceiveMessage - handle the messages incoming from the websocket
func (h *Hub) HandleReceiveMessage(client Client, messageType int, payload []byte) *Hub {
	m := Message{}

	// test data: '{"type": 1, "body": "foobar", "senderID": 110, "clientIDS": [123, 456, 789]}'
	err := json.Unmarshal(payload, &m)

	if err != nil {
		fmt.Println("Unrecognised message")
		return h
	}

	fmt.Printf("Valid payload :)\n MsgType: %v\n Body: %v\n SenderID: %v\n ClientIDS: %+v\n", m.MsgType, m.Body, m.SenderID, m.ClientIDS)
	return h
}
