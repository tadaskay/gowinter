package event

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	tables := []struct {
		event interface{}
		msg   string
	}{
		{StartEvent{Name: "john"}, "START john"},
		{WalkEvent{5, 1}, "WALK 5 1"},
		{ShootEvent{2, 2}, "SHOOT 2 2"},
		{MissEvent{}, "MISS"},
		{BoomEvent{"john", "night-king"}, "BOOM john night-king"},
		{VictoryEvent{"Petras"}, "VICTORY Petras"},
		{DefeatEvent{}, "DEFEAT"},
	}

	for _, table := range tables {
		msg, _ := Marshal(table.event)

		if table.msg != msg {
			t.Errorf("Expected message: '%v', but got: '%v'", table.msg, msg)
		}
	}
}
