package game

import (
	"fmt"
	"github.com/tadaskay/gowinter/game/board"
	"github.com/tadaskay/gowinter/game/network"
	"github.com/tadaskay/gowinter/game/zombie"
)

type Board struct {
	size   board.Bounds
	zombie zombie.Zombie
}

type Player string

type State int

const (
	Pregame State = iota + 1
	Started
	Finished
)

type Session struct {
	board  board.Bounds
	zombie *zombie.Zombie
	player Player
	client *network.GameClient
	state  State
	End    chan bool
}

func NewSession(sizeX, sizeY int, client *network.GameClient) Session {
	session := Session{
		board:  board.Bounds{sizeX, sizeY},
		zombie: zombie.New("night-king"),
		client: client,
		state:  Started,
		End:    make(chan bool),
	}
	go session.gameLoop()
	return session
}

func (session *Session) gameLoop() {
	for {
		session.processInput()
		session.update()
	}
}

func (session *Session) processInput() {
	select {
	case clientEvent := <-session.client.Events:
		fmt.Println("Received: ", clientEvent)
	default:
	}
}

func (session *Session) update() {
	if session.state == Started {
		if session.zombie.State == zombie.Initial {
			session.zombie.Spawn(session.board)
		} else {
			session.zombie.Update()
		}

		if session.zombie.IsSouthReached() {
			session.state = Finished
			session.End <- true
		}
	}
}
