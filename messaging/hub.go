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
	fmt.Println("Made it this far...", h.subscriptions)
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
	// clientSubs := h.GetSubscriptions(client)

	// if len(clientSubs) > 0 {
	// 	// client is subscribed
	// 	return h
	// }

	s := newSubscription(topic, client)
	h.subscriptions = append(h.subscriptions, *s)

	fmt.Printf("Client %v has been subscribed to the %s topic \n", client.userID, topic)
	return h
}

func (h *Hub) Publish(message []byte, excludeClient *Client) {
	subscriptions := h.GetSubscriptions(nil)
	//fmt.Printf("subs: %+v", subscriptions)
	//if len(subscriptions) == len(h.clients) {
	//fmt.Println("Not Balls")
	for _, sub := range subscriptions {
		//fmt.Printf("Sending to client id %v message is %s", sub.client.userID, message)
		if sub.client != nil {
			fmt.Println("Here?", sub)
			//sub.client.connection.WriteMessage(1, message)
			err := sub.client.connection.WriteMessage(1, message)

			if err != nil {
				log.Println(err)
				return
			}
		}

	}
	// } else {
	// 	fmt.Println("Balls")
	// }

}

// HandleReceiveMessage - handle the messages incoming from the websocket
func (h *Hub) HandleReceiveMessage(client Client, messageType int, payload []byte) *Hub {
	m := Message{}

	// test data: '{"type": 1, "body": "foobar", "senderID": 110, "clientIDS": [123, 456, 789]}'
	err := json.Unmarshal(payload, &m)

	if err != nil {
		fmt.Println("Unrecognised message format")
		return h
	}

	fmt.Printf("Valid payload :)\n"+
		"MsgType: %v\n Body: %v\n SenderID: %v\n ClientIDS: %+v\n",
		m.MsgType,
		string(m.Body),
		m.SenderID,
		m.ClientIDS)

	// switch for identity, list, relay
	h.Publish(m.Body, nil)
	return h
}
