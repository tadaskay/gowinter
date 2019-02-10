package zombie

import (
	"fmt"
	"math/rand"
	"time"
	"winter-is-coming/game/board"
)

type Zombie struct {
	Name string
	Pos  board.Position

	bounds     board.Bounds
	moveTicker *time.Ticker
}

func (z *Zombie) OnStart(bounds board.Bounds) {
	rand.Seed(time.Now().UnixNano())
	z.Pos = board.Position{X: rand.Intn(bounds.X)}
	z.bounds = bounds
	z.moveTicker = time.NewTicker(2 * time.Second)
}

func (z *Zombie) Update() {
	select {
	case _ = <-z.moveTicker.C:
		z.move()
	default:
	}
}

// Moves to a new position south:
// bottom-left / bottom / bottom-right at random
func (z *Zombie) move() {
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
		fmt.Println("WALL REACHED", z.Pos)
		z.moveTicker.Stop()
		return
	} else {
		fmt.Println("MOVE", z.Pos)
	}
}

func (z *Zombie) IsSouthReached() bool {
	return z.Pos.Y == z.bounds.Y
}
