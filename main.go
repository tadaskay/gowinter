package main

import (
	"github.com/tadaskay/gowinter/game"
	"github.com/tadaskay/gowinter/network"
)

var PORT = 52000

func main() {
	noopKillChannel, noopDoneChannel := make(chan bool), make(chan bool)
	interruptableMain(noopKillChannel, noopDoneChannel)
}

// For testing purposes, it's interruptible from outside
func interruptableMain(kill chan bool, done chan bool) {
	client := network.NewGameClient(PORT)
	session := game.NewSession(10, 30, client)
	for {
		select {
		case _ = <-kill:
			done <- true
			return
		case _ = <-session.Complete:
			done <- true
			return
		default:
		}
	}
}
