package messaging

import "encoding/json"

type messageType int

const (
	// Identity - Used to return the client's userID back to them
	Identity messageType = iota

	// List - Used to return a list of all connected userID's (excluding the requesting client)
	List

	// Relay - Used to relay a given message body to selected receivers in the message
	Relay
)

// Message - used by the Hub to send to clients
type Message struct {
	MsgType   messageType     `json:"type"`
	Body      json.RawMessage `json:"body"`
	ClientIDS []uint64        `json:"clientIDS"`
}
