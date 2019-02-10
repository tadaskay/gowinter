package game

import (
	"fmt"
	"github.com/tadaskay/gowinter/board"
	"github.com/tadaskay/gowinter/event"
	"github.com/tadaskay/gowinter/network"
	"github.com/tadaskay/gowinter/zombie"
)

type Player string

type State int

const (
	Pregame State = iota + 1
	Started
	Finished
)

type Session struct {
	board        board.Bounds
	zombie       *zombie.Zombie
	player       Player
	client       *network.GameClient
	serverEvents chan interface{}
	state        State
	Complete     chan bool
}

func NewSession(sizeX, sizeY int, client *network.GameClient) Session {
	serverEvents := make(chan interface{}, 5)
	session := Session{
		board:        board.Bounds{sizeX, sizeY},
		zombie:       zombie.New("night-king", serverEvents),
		client:       client,
		serverEvents: serverEvents,
		state:        Pregame,
		Complete:     make(chan bool),
	}
	go session.gameLoop()
	return session
}

func (session *Session) gameLoop() {
	defer func() {
		session.Complete <- true
	}()
	for {
		session.processEvents()
		session.update()

		if session.state == Finished {
			if len(session.client.Sent) == 0 && len(session.serverEvents) == 0 {
				return
			}
		}
	}
}

func (session *Session) processEvents() {
	session.handleClientEvents()
	session.sendServerEvents()
}

func (session *Session) handleClientEvents() {
	select {
	case evt := <-session.client.Received:
		msg, _ := event.Marshal(evt)
		fmt.Println("Client event:", msg)
		switch clientEvent := evt.(type) {
		case event.StartEvent:
			session.start(clientEvent.Name)
		case event.ShootEvent:
			session.shotsFired(string(session.player), clientEvent.X, clientEvent.Y)
		}
	default:
	}
}

func (session *Session) sendServerEvents() {
	select {
	case evt := <-session.serverEvents:
		msg, _ := event.Marshal(evt)
		fmt.Println("Server event:", msg)
		session.client.Sent <- evt
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

func (session *Session) shotsFired(playerName string, x, y int) {
	session.zombie.HandleShot(playerName, board.Position{x, y})
}

func (session *Session) determineIfGameFinished() {
	if session.zombie.IsSouthReached() {
		session.serverEvents <- event.EndEvent{Victory: false}
		session.state = Finished
	}
	if session.zombie.IsDead() {
		session.serverEvents <- event.EndEvent{Victory: true}
		session.state = Finished
	}
}
