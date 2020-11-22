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

// NewHub creates a new hub
func NewHub(clients []Client, subs []Subscription) *Hub {
	return &Hub{clients, subs}
}

// Subscription - allows client to subscribe to the topic
type Subscription struct {
	topic  string
	client *Client
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

// Subscribe - creates new subcription to topic
func (h *Hub) Subscribe(client *Client) *Hub {
	s := newSubscription(topic, client)
	h.subscriptions = append(h.subscriptions, *s)

	fmt.Printf("Client %v has been subscribed to the %s topic \n", client.userID, topic)
	return h
}

// getSubscriptions - returns a slice of all Subscriptions in the Hub, omits sender if client is nil
func (h *Hub) getSubscriptions(client *Client) []Subscription {
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

// getRequestedSubscriptions - returns all subscriptions from a slice of userIDs
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

// publishToSender - publishes a message back to the original sender only
func (h *Hub) publishToSender(message []byte, client *Client) (isOK bool) {
	err := client.connection.WriteMessage(1, message)

	if err != nil {
		log.Println(err)
		return false
	}

	fmt.Printf("Sending to client id %v with message %s \n", client.userID, message)
	return true
}

// publishToReceivers - publishes a message to each subscription that exists based on the subs that are passed in
func (h *Hub) publishToReceivers(message []byte, subs []Subscription) (isOK bool) {
	for _, s := range subs {
		if s.client != nil {
			err := s.client.connection.WriteMessage(1, message)

			if err != nil {
				log.Println(err)
				return false
			}

			fmt.Printf("Sending to client id %v with message %s \n", s.client.userID, message)
		}
	}
	return true
}

/*
	handleIdentity - handles the Identity message type (Used to return the senders userID back to them)
	Test data: '{"type": 0 }'
*/
func (h *Hub) handleIdentity(client *Client) {
	/* We could simply get the sender userID from the client
	but it should really be checked via the subscriptions */
	subs := h.getSubscriptions(client)
	var id string

	if len(subs) == 1 && subs[0].client.userID == client.userID {
		id = strconv.FormatUint(subs[0].client.userID, 10)
	} else {
		fmt.Println("Something went wrong - this client is not subscribed")
	}

	payload := "(Identity) Current userID: " + id
	msg, isOK := SerialiseString(payload)

	if isOK {
		isOK := h.publishToSender(msg, client)
		if !isOK {
			fmt.Println("Failed to publish to sender")
		}
	}
}

/*
	handleList - handles the List message type (Used to return a list of all connected userID's (excluding the requesting client))
	Test data: '{"type": 1 }'
*/
func (h *Hub) handleList(client *Client) {
	subs := h.getSubscriptions(nil)
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
		payload = "(List) It's dangerous to go alone! Take this. http://localhost:8888 "
	}

	msg, isOK := SerialiseString(payload)

	if isOK {
		isOK := h.publishToSender(msg, client)
		if !isOK {
			fmt.Println("Failed to publish to sender")
		}
	}
}

/*
	handleRelay - handles the Relay message type (Used to relay a given message body to selected receivers in the message)
	Test data (where 1, 2, 3 has to be existing userIDs, use a List call above first :) ): '{"type": 2, "body": "foobar", "clientIDS": [1, 2, 3]}'
*/
func (h *Hub) handleRelay(client *Client, message *Message) {
	sort.SliceStable(message.ClientIDS, func(i, j int) bool {
		return message.ClientIDS[i] < message.ClientIDS[j]
	})

	subs := h.getRequestedSubscriptions(message.ClientIDS)
	body, isOk := DeserialiseString(message.Body)

	if isOk {
		var payload string
		okSubs, okBody, retMsg := ValidateRequest(subs, body)

		if okSubs && okBody {
			if len(subs) > 0 {
				payload = "(Relay) - " + body
				msg, isOK := SerialiseString(payload)

				if isOK {
					isOK := h.publishToReceivers(msg, subs)
					if !isOK {
						fmt.Println("Failed to publish")
					}
				}
			} else {
				payload = "There are no clients that match that/those userID/s"
				msg, isOK := SerialiseString(payload)

				if isOK {
					isOK := h.publishToSender(msg, client)
					if !isOK {
						fmt.Println("Failed to publish to sender")
					}
				}

			}
		} else if !okSubs {
			msg, isOK := SerialiseString(retMsg)

			if isOK {
				isOK := h.publishToSender(msg, client)
				if !isOK {
					fmt.Println("Failed to publish to sender")
				}
			}
		} else if !okBody {
			msg, isOK := SerialiseString(retMsg)

			if isOK {
				isOK := h.publishToSender(msg, client)
				if !isOK {
					fmt.Println("Failed to publish to sender")
				}
			}
		}
	}

}

/*
	handleDefault - handles the default case when the message type is not recognised
	Test data: '{"type": 3 }'
*/
func (h *Hub) handleDefault(client *Client) {
	payload := "Unrecognised message type, please use ints only: " +
		"0: Identity, 1: List, 2: Relay"
	msg, isOK := SerialiseString(payload)

	if isOK {
		fmt.Println(payload)
		isOK := h.publishToSender(msg, client)
		if !isOK {
			fmt.Println("Failed to publish to sender")
		}
	}
}

// HandleReceiveMessage - handles the messages incoming from the websocket
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
