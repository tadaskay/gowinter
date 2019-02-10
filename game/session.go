package game

import (
	"fmt"
	"github.com/tadaskay/gowinter/game/board"
	"github.com/tadaskay/gowinter/game/event"
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
		state:  Pregame,
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
	case evt := <-session.client.Events:
		msg, _ := event.Marshal(evt)
		fmt.Println("Client event:", msg)
		switch clientEvent := evt.(type) {
		case event.StartEvent:
			session.start(clientEvent.Name)
		}
	default:
	}
}

func (session *Session) update() {
	if session.state == Started {
		session.zombie.Update()
		session.determineIfGameFinished()
	}
}

func (session *Session) start(playerName string) {
	session.player = Player(playerName)
	session.state = Started
	session.zombie.Spawn(session.board)
}

func (session *Session) determineIfGameFinished() {
	if session.zombie.IsSouthReached() {
		session.state = Finished
		session.End <- true
	}
}
