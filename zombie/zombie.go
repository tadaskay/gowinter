package zombie

import (
	"github.com/tadaskay/gowinter/board"
	"github.com/tadaskay/gowinter/event"
	"math/rand"
	"time"
)

type Zombie struct {
	Name  string
	pos   board.Position
	state state

	bounds       board.Bounds
	moveTicker   *time.Ticker
	serverEvents chan<- interface{}
}

type state int

const (
	initial state = iota
	moving
	dead
)

func New(name string, serverEvents chan interface{}) *Zombie {
	return &Zombie{
		Name:         name,
		state:        initial,
		serverEvents: serverEvents,
	}
}

func (z *Zombie) Spawn(bounds board.Bounds) {
	rand.Seed(time.Now().UnixNano())
	z.pos = board.Position{X: rand.Intn(bounds.X)}
	z.bounds = bounds
	z.moveTicker = time.NewTicker(2 * time.Second)
	z.state = moving
}

func (z *Zombie) HandleShot(playerName string, pos board.Position) {
	if z.pos == pos {
		z.state = dead
		z.serverEvents <- event.BoomEvent{Character: playerName, Target: z.Name}
	} else {
		z.serverEvents <- event.MissEvent{}
	}
}

func (z *Zombie) Update() {
	if z.state == moving {
		z.move()
	}
}

// Moves to a new position south:
// bottom-left / bottom / bottom-right at random
func (z *Zombie) move() {
	select {
	case _ = <-z.moveTicker.C:
	default:
		return
	}

	deltaX := rand.Intn(3) - 1
	newX := z.pos.X + deltaX
	if newX < 0 {
		newX = 0
	} else if newX > z.bounds.X {
		newX = z.bounds.X
	}

	z.pos.X = newX
	z.pos.Y += 1

	if z.IsSouthReached() {
		z.moveTicker.Stop()
		return
	} else {
		z.serverEvents <- event.WalkEvent{z.pos.X, z.pos.Y}
	}
}

func (z *Zombie) IsSouthReached() bool {
	return z.pos.Y == z.bounds.Y
}

func (z *Zombie) IsDead() bool {
	return z.state == dead
}
