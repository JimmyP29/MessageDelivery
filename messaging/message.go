package messaging

type messageType int

const (
	identity messageType = iota
	list
	relay
)

type message struct {
	msgType messageType
	body    string
	client  *client
}
