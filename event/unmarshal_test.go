package event

import (
	"errors"
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	tables := []struct {
		msg   string
		event interface{}
		err   error
	}{
		{"", nil, e("unrecognized event id [], input: []")},

		{"START john", StartEvent{Name: "john"}, nil},
		{"START", nil, e("1 argument(s) required, input: [START]")},

		{"WALK 5 1", WalkEvent{5, 1}, nil},
		{"WALK A 1", nil, e("cannot read int argument [A], input: [WALK A 1]")},

		{"SHOOT 2 2", ShootEvent{2, 2}, nil},
		{"SHOOT 2 B", nil, e("cannot read int argument [B], input: [SHOOT 2 B]")},
		{"MISS", MissEvent{}, nil},
		{"BOOM john night-king", BoomEvent{"john", "night-king"}, nil},
		{"VICTORY Petras", VictoryEvent{"Petras"}, nil},
		{"DEFEAT", DefeatEvent{}, nil},
	}

	for _, table := range tables {
		event, err := Unmarshal(table.msg)

		if fmt.Sprint(table.err) != fmt.Sprint(err) {
			t.Errorf("Expected error: '%v', but got: '%v'", table.err, err)
		}

		if table.event != event {
			t.Errorf("Expected unmarshalled event: '%v', but got: '%v'", table.event, event)
		}
	}
}

func e(msg string) error {
	return errors.New(msg)
}
