package event

import "reflect"

type id string

const (
	START id = "START"
	WALK  id = "WALK"
	SHOOT id = "SHOOT"
	MISS  id = "MISS"
	BOOM  id = "BOOM"
	END   id = "END"
)

var Types = map[id]reflect.Type{
	START: reflect.TypeOf(StartEvent{}),
	WALK:  reflect.TypeOf(WalkEvent{}),
	SHOOT: reflect.TypeOf(ShootEvent{}),
	MISS:  reflect.TypeOf(MissEvent{}),
	BOOM:  reflect.TypeOf(BoomEvent{}),
	END:   reflect.TypeOf(EndEvent{}),
}

type StartEvent struct {
	Name string
}

type WalkEvent struct {
	X, Y int
}

type ShootEvent struct {
	X, Y int
}

type MissEvent struct {
}

type BoomEvent struct {
	Character string
	Target    string
}

type EndEvent struct {
	Victory bool
}
