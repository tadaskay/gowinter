package main

import (
	"github.com/tadaskay/gowinter/game"
	"github.com/tadaskay/gowinter/network"
)

var PORT = 52000

func main() {
	client := network.NewGameClient(PORT)
	session := game.NewSession(10, 30, client)
	<-session.Complete
}
