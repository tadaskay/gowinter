package game

import (
	"fmt"
	"math/rand"
	"time"
)

type Bounds struct {
	x, y int
}

type Position struct {
	x, y int
}

type Board struct {
	size   Bounds
	zombie Zombie
}

type Zombie struct {
	name string
	pos  Position
}

func (z *Zombie) spawn(bounds Bounds) {
	z.pos = Position{rand.Intn(bounds.x), 0}
}

func (z *Zombie) moveSouth(bounds Bounds) {
	deltaX := rand.Intn(3) - 1
	newX := z.pos.x + deltaX
	if newX < 0 {
		newX = 0
	} else if newX > bounds.x {
		newX = bounds.x
	}
	z.pos.x = newX
	z.pos.y += 1
}

func (board *Board) spawn(z Zombie) {
	board.zombie = z
	board.zombie.spawn(Bounds{board.size.x, 0})
}

func (board *Board) ZombiePosition() {
	fmt.Println(board.zombie.pos)
}

func (board *Board) MoveZombie() {
	moveTick := time.Tick(2 * time.Second)
	for range moveTick {
		board.zombie.moveSouth(board.size)
		board.ZombiePosition()
	}
}

type Player string

type Session struct {
	Board  Board
	player Player
	End    chan bool
}

func (session *Session) Start() {
	go session.Board.MoveZombie()
}

func NewSession(sizeX, sizeY int, playerName string) Session {
	rand.Seed(time.Now().UnixNano())
	board := Board{size: Bounds{sizeX, sizeY}}
	zombie := Zombie{name: "night-king"}
	board.spawn(zombie)
	session := Session{
		board,
		Player(playerName),
		make(chan bool),
	}
	session.Start()
	return session
}
