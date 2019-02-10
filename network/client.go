package network

import (
	"fmt"
	"github.com/tadaskay/gowinter/event"
	"net"
	"os"
	"strconv"
	"strings"
)

type GameClient struct {
	socket net.Conn
	Events chan interface{}
}

func NewGameClient(port int) *GameClient {
	conn := waitForClientConnection(port)
	client := &GameClient{
		socket: conn,
		Events: make(chan interface{}),
	}
	go client.receive()
	return client
}

func waitForClientConnection(port int) net.Conn {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error listening for connections", err)
		os.Exit(1)
	}
	conn, err := listener.Accept()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error connecting: ", err)
		os.Exit(1)
	}
	return conn
}

func (client *GameClient) receive() {
	defer func() {
		_ = client.socket.Close()
	}()
	for {
		buf := make([]byte, 4096)
		n, err := client.socket.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			break
		}

		received := n > 0
		if !received {
			continue
		}

		message := strings.TrimRight(string(buf[:n]), "\r\n")
		gameEvent, err2 := event.Unmarshal(message)
		if err2 != nil {
			fmt.Println("Invalid message from client:", err2)
			continue
		}
		client.Events <- gameEvent
	}
}
