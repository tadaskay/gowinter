package event

import (
	"reflect"
)

type id string

const (
	START id = "START"
	WALK  id = "WALK"
	SHOOT id = "SHOOT"
	MISS  id = "MISS"
	BOOM  id = "BOOM"
	END   id = "END"
)

var idToType = map[id]reflect.Type{
	START: reflect.TypeOf(StartEvent{}),
	WALK:  reflect.TypeOf(WalkEvent{}),
	SHOOT: reflect.TypeOf(ShootEvent{}),
	MISS:  reflect.TypeOf(MissEvent{}),
	BOOM:  reflect.TypeOf(BoomEvent{}),
	END:   reflect.TypeOf(EndEvent{}),
}

var typeToIdCache = make(map[reflect.Type]id)

func typeToId(t reflect.Type) (result id, found bool) {
	if id, found := typeToIdCache[t]; found {
		return id, true
	}

	if len(typeToIdCache) == 0 {
		for k, v := range idToType {
			typeToIdCache[v] = k
		}
		var id, found = typeToIdCache[t]
		return id, found
	}

	return "", false
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
