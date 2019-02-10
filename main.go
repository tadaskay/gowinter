package main

import (
	"fmt"
	"github.com/tadaskay/gowinter/game"
	"github.com/tadaskay/gowinter/network"
	"net"
	"os"
	"strconv"
)

var PORT = 52000

func main() {
	conn := waitForClientConnection()
	client := network.NewGameClient(conn)

	session := game.NewSession(10, 30, client)
	<-session.End
	fmt.Println("Game ended")
}

func waitForClientConnection() net.Conn {
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
	return conn
}
