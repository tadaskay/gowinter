package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"winter-is-coming/game"
)

type Client struct {
	socket net.Conn
	data   chan []byte
}

var PORT int = 52000

func receive(client *Client) {
	for {
		message := make([]byte, 4096)
		n, err := client.socket.Read(message)
		if err != nil {
			fmt.Println("Error reading from client: ", err)
			_ = client.socket.Close()
			break
		}
		if n > 0 {
			client.data <- message[:n]
		}
	}
}

func main() {
	var client *Client
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error listening for connections", err)
		os.Exit(1)
	}
	conn, err := listener.Accept()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "Error connecting: ", err)
		os.Exit(1)
	}
	client = &Client{socket: conn, data: make(chan []byte)}
	go receive(client)
	message := <-client.data
	fmt.Println(string(message))

	session := game.NewSession(10, 30, "john")
	session.Start()
	<-session.End
}
