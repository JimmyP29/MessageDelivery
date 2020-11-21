package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var topic = "Multiplay"

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

// GetSubscriptions - returns a slice of all Subscriptions in the Hub, omits sender if client is nil
func (h *Hub) GetSubscriptions(client *Client) []Subscription {
	var subs []Subscription

	for _, s := range h.subscriptions {
		if s.client != nil {
			if client != nil {
				if s.client.userID == client.userID {
					subs = append(subs, s)
				}
			} else {
				subs = append(subs, s)
			}
		}
	}

	return subs
}

// GetRequestedSubscriptions - returns all subscriptions from a slice of userIDs
func (h *Hub) getRequestedSubscriptions(ids []uint64) []Subscription {
	var subs []Subscription

	for _, s := range h.subscriptions {
		if s.client != nil {
			inBoth := false
			for _, i := range ids {
				if i == s.client.userID {
					inBoth = true
					break
				}
			}

			if inBoth {
				subs = append(subs, s)
			}
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

// publishToSender - publishes a message back to the original sender only
func (h *Hub) publishToSender(message []byte, client *Client) {
	err := client.connection.WriteMessage(1, message)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Sending to client id %v with message %s \n", client.userID, message)
}

func (h *Hub) publishToReceivers(message []byte, subs []Subscription) {
	for _, s := range subs {
		if s.client != nil {
			err := s.client.connection.WriteMessage(1, message)

			if err != nil {
				log.Println(err)
				return
			}

			fmt.Printf("Sending to client id %v with message %s \n", s.client.userID, message)
		}
	}
}

/*
	handleIdentity - handles the Identity message type (Used to return the senders userID back to them)
	Test data: '{"type": 0 }'
*/
func (h *Hub) handleIdentity(client *Client) {
	/* We could simply get the sender userID from the client
	but it should really be checked via the subscriptions */
	subs := h.GetSubscriptions(client)
	var id string

	if len(subs) == 1 && subs[0].client.userID == client.userID {
		id = strconv.FormatUint(subs[0].client.userID, 10)
	} else {
		fmt.Println("Something went wrong - this client is not subscribed")
	}

	payload := "(Identity) Current userID: " + id
	msg, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
		return
	}

	h.publishToSender(msg, client)
}

/*
	handleList - handles the List message type (Used to return a list of all connected userID's (excluding the requesting client))
	Test data: '{"type": 1 }'
*/
func (h *Hub) handleList(client *Client) {
	subs := h.GetSubscriptions(nil)
	var returnIDs []string

	for _, s := range subs {
		if s.client != nil {
			if s.client.userID != client.userID {
				returnIDs = append(returnIDs, strconv.FormatUint(s.client.userID, 10))
			}
		}

	}

	var payload string
	if len(returnIDs) >= 1 {
		payload = "(List) Other current userIDs: " + strings.Join(returnIDs, ", ")
	} else {
		payload = "(List) You are all alone!"
	}

	msg, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
		return
	}

	h.publishToSender(msg, client)
}

/*
	handleRelay - handles the Relay message type (Used to relay a given message body to selected receivers in the message)
	Test data (where 1, 2, 3 has to be existing userIDs, use a List call above first :) ): '{"type": 2, "body": "foobar", "clientIDS": [1, 2, 3]}'
*/
func (h *Hub) handleRelay(client *Client, message *Message) {

	fmt.Printf("before: %+v", message.ClientIDS)
	sort.SliceStable(message.ClientIDS, func(i, j int) bool {
		return message.ClientIDS[i] < message.ClientIDS[j]
	})
	subs := h.getRequestedSubscriptions(message.ClientIDS)

	var payload string
	if len(subs) > 0 {
		bytes, err := json.Marshal(message.Body)
		if err != nil {
			log.Println(err)
			return
		}

		var body string
		err = json.Unmarshal(bytes, &body)
		if err != nil {
			log.Println(err)
			return
		}

		payload = "(Relay) - " + body
		msg, err := json.Marshal(payload)

		if err != nil {
			log.Println(err)
			return
		}

		h.publishToReceivers(msg, subs)
	} else {
		payload = "There are no clients that match that/those userID/s"
		fmt.Println(payload)

		msg, err := json.Marshal(payload)

		if err != nil {
			log.Println(err)
			return
		}

		h.publishToSender(msg, client)
	}

}

/*
	handleDefault - handles the default case when the message type is not recognised
	Test data: '{"type": 3 }'
*/
func (h *Hub) handleDefault(client *Client) {
	payload := "Unrecognised message type, please use ints only: " +
		"0: Identity, 1: List, 2: Relay"

	msg, err := json.Marshal(payload)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(payload)
	h.publishToSender(msg, client)
}

// HandleReceiveMessage - handle the messages incoming from the websocket
func (h *Hub) HandleReceiveMessage(client Client, payload []byte) *Hub {
	m := Message{}
	err := json.Unmarshal(payload, &m)

	if err != nil {
		fmt.Println("Unrecognised message format")
		return h
	}

	if &client != nil {
		switch m.MsgType {
		case Identity:
			h.handleIdentity(&client)
			break
		case List:
			h.handleList(&client)
			break
		case Relay:
			h.handleRelay(&client, &m)
			break
		default:
			h.handleDefault(&client)
			break
		}
	} else {
		fmt.Println("Invalid client")
	}

	return h
}
