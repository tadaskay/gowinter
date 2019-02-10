package game

import (
	"winter-is-coming/game/board"
	"winter-is-coming/game/zombie"
)

type Board struct {
	size   board.Bounds
	zombie zombie.Zombie
}

type Player string

type Session struct {
	board  board.Bounds
	zombie zombie.Zombie
	player Player
	End    chan bool
}

func NewSession(sizeX, sizeY int, playerName string) Session {
	session := Session{
		board:  board.Bounds{sizeX, sizeY},
		zombie: zombie.Zombie{Name: "night-king"},
		player: Player(playerName),
		End:    make(chan bool),
	}
	return session
}

func (session *Session) Start() {
	session.zombie.OnStart(session.board)
	go session.gameLoop()
}

func (session *Session) gameLoop() {
	for {
		session.zombie.Update()
		if session.zombie.IsSouthReached() {
			session.End <- true
			break
		}
	}
}
