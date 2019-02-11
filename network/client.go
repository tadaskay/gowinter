package network

import (
	"fmt"
	"github.com/tadaskay/gowinter/event"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type GameClient struct {
	socket   net.Conn
	Received chan interface{}
	Sent     chan interface{}
	fatal    chan error
}

func NewGameClient(port int) *GameClient {
	fmt.Println("Waiting for client connection on port:", port)
	conn := waitForClientConnection(port)
	client := &GameClient{
		socket:   conn,
		Received: make(chan interface{}, 5),
		Sent:     make(chan interface{}, 5),
	}
	go client.send()
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

func (client *GameClient) send() {
	defer client.socket.Close()
	for {
		select {
		case evt, ok := <-client.Sent:
			if !ok {
				return
			}
			msg, err := event.Marshal(evt)
			if err != nil {
				fmt.Println("Error marshalling message to client:", msg)
				continue
			}
			_, err2 := io.WriteString(client.socket, msg+"\r\n")
			if err2 != nil {
				fmt.Println("Fatal error occurred when writing to client", err2)
				return
			}
		}
	}
}

func (client *GameClient) receive() {
	defer client.socket.Close()
	for {
		buf := make([]byte, 4096)
		n, err := client.socket.Read(buf)
		if err != nil {
			fmt.Println("Fatal error occurred when reading with client", err)
			return
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

		client.Received <- gameEvent
	}
}
