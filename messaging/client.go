package messaging

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	userID   uint64
	messages chan<- message
}

func (c *client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		messages := strings.TrimSpace(args[0])

		fmt.Printf("messages: %v", messages)
		switch messages {
		case "/identity":
			c.messages <- message{
				msgType: identity,
				client:  c,
			}
		case "/list":
			c.messages <- message{
				msgType: list,
				client:  c,
			}
		case "/relay":
			c.messages <- message{
				msgType: relay,
				client:  c,
			}
		default:
			c.err(fmt.Errorf("Unknown command: %s", messages))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("Error: " + err.Error() + "\n"))
}

func (c *client) WriteMessage(msg *message) {
	c.conn.Write([]byte("testing" + msg.body))
}
