package network

import (
	"fmt"
	"github.com/tadaskay/gowinter/game/event"
	"net"
	"strings"
)

type GameClient struct {
	socket net.Conn
	Events chan interface{}
}

func NewClient(conn net.Conn) *GameClient {
	client := &GameClient{
		socket: conn,
		Events: make(chan interface{}),
	}
	return client
}

func (client *GameClient) StartReceiving() {
	go func() {
		for {
			buf := make([]byte, 4096)
			n, err := client.socket.Read(buf)
			if err != nil {
				fmt.Println("Error reading from client:", err)
				_ = client.socket.Close()
				break
			}
			if n > 0 {
				message := strings.TrimRight(string(buf[:n]), "\r\n")
				gameEvent, err := event.Unmarshal(message)
				if err != nil {
					fmt.Println("Invalid message from client:", err)
				}
				client.Events <- gameEvent
			}
		}
	}()
}
