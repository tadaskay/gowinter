package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"winter-is-coming/game"
	"winter-is-coming/game/network"
)

var PORT int = 52000

func main() {
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
	client := network.NewClient(conn)
	client.StartReceiving()

	session := game.NewSession(10, 30, client)
	<-session.End
	fmt.Println("Game ended")
}
