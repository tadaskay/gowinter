package zombie

import (
	"fmt"
	"github.com/tadaskay/gowinter/board"
	"math/rand"
	"time"
)

type Zombie struct {
	Name  string
	Pos   board.Position
	state State

	bounds     board.Bounds
	moveTicker *time.Ticker
}

type State int

const (
	Initial State = iota
	Moving
	Dead
)

func New(name string) *Zombie {
	return &Zombie{
		Name:  name,
		state: Initial,
	}
}

func (z *Zombie) Spawn(bounds board.Bounds) {
	rand.Seed(time.Now().UnixNano())
	z.Pos = board.Position{X: rand.Intn(bounds.X)}
	z.bounds = bounds
	z.moveTicker = time.NewTicker(2 * time.Second)
	z.state = Moving
}

func (z *Zombie) HandleShot(pos board.Position) {
	if z.Pos == pos {
		z.state = Dead
	}
}

func (z *Zombie) Update() {
	if z.state == Moving {
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
	newX := z.Pos.X + deltaX
	if newX < 0 {
		newX = 0
	} else if newX > z.bounds.X {
		newX = z.bounds.X
	}

	z.Pos.X = newX
	z.Pos.Y += 1

	if z.IsSouthReached() {
		z.moveTicker.Stop()
		return
	} else {
		fmt.Println("MOVE", z.Pos)
	}
}

func (z *Zombie) IsSouthReached() bool {
	return z.Pos.Y == z.bounds.Y
}

func (z *Zombie) IsDead() bool {
	return z.state == Dead
}
