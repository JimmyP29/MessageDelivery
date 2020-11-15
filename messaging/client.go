package messaging

import "net"

type client struct {
	conn     net.Conn
	userID   uint64
	messages chan<- message
}
