package messaging

import "encoding/json"

type messageType int

const (
	identity messageType = iota
	list
	relay
)

// Message - used by the Hub to send to clients
type Message struct {
	MsgType   messageType     `json:"type"`
	Body      json.RawMessage `json:"body"`
	SenderID  uint64          `json:"senderID"`
	ClientIDS []uint64        `json:"clientIDS"`
}
