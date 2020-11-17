package messaging

import (
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
	return uint64(rand.Intn(255-1) + 1)
}

// AddClient - appends new client to client slice
func (h *Hub) AddClient(client Client) *Hub {
	h.clients = append(h.clients, client)

	fmt.Printf("Client %v has been added top the hub \n", client.userID)
	fmt.Printf("hub: %+v \n", h.clients)
	return h
}
