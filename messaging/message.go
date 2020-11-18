package messaging

type messageType int

const (
	identity messageType = iota
	list
	relay
)

// Message - used by the Hub to send to clients
type Message struct {
	MsgType   messageType `json:"type"`
	Body      string      `json:"body"`
	SenderID  uint64      `json:"senderID"`
	ClientIDS []uint64    `json:"clientIDS"`
	Topic     string
}
