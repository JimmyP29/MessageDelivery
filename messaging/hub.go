package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var topic = "TechTest"

// Hub - used to control message flow to clients
type Hub struct {
	clients       []Client
	subscriptions []Subscription
}

// Subscription - allows client to subscribe to the topic
type Subscription struct {
	topic  string
	client *Client
}

// NewHub creates a new hub
func NewHub(clients []Client, subs []Subscription) *Hub {
	return &Hub{clients, subs}
}

func newSubscription(topic string, client *Client) *Subscription {
	return &Subscription{topic, client}
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
	return h
}

// GetSubscriptions - returns a slice of all Subscriptions in the Hub
func (h *Hub) GetSubscriptions(client *Client) []Subscription {
	var subs []Subscription

	for _, sub := range h.subscriptions {
		if client != nil {
			if sub.client.userID == client.userID {
				subs = append(subs, sub)
			}
		} else {
			subs = append(subs, sub)
		}
	}

	return subs
}

// Subscribe - creates new subcription to topic
func (h *Hub) Subscribe(client *Client) *Hub {
	s := newSubscription(topic, client)
	h.subscriptions = append(h.subscriptions, *s)

	fmt.Printf("Client %v has been subscribed to the %s topic \n", client.userID, topic)
	return h
}

// Publish - broadcasts a message on the websocket
func (h *Hub) Publish(message []byte, excludeClient *Client) {
	subscriptions := h.GetSubscriptions(nil)

	for _, sub := range subscriptions {
		if sub.client != nil {
			err := sub.client.connection.WriteMessage(1, message)

			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Sending to client id %v with message %s \n", sub.client.userID, message)
		}
	}
}

// HandleReceiveMessage - handle the messages incoming from the websocket
func (h *Hub) HandleReceiveMessage(client Client, payload []byte) *Hub {
	m := Message{}

	// test data: '{"type": 1, "body": "foobar", "senderID": 110, "clientIDS": [123, 456, 789]}'
	err := json.Unmarshal(payload, &m)

	if err != nil {
		fmt.Println("Unrecognised message format")
		return h
	}

	// fmt.Printf("Valid payload :)\n"+
	// 	"MsgType: %v\n Body: %v\n SenderID: %v\n ClientIDS: %+v\n",
	// 	m.MsgType,
	// 	string(m.Body),
	// 	m.SenderID,
	// 	m.ClientIDS)

	// switch for identity, list, relay
	switch m.MsgType {
	case Identity:
		fmt.Printf("Identity: %v \n", Identity)
		break
	case List:
		fmt.Printf("List: %v \n", List)
		break
	case Relay:
		fmt.Printf("Relay: %v \n", Relay)
		break
	default:
		fmt.Println("Unrecognised message type, please use: \n" +
			"0: Identity\n" + "1: List\n" + "2: Relay\n")
		break
	}
	h.Publish(m.Body, nil)
	return h
}
