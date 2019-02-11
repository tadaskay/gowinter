package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		msg   string
		event interface{}
		err   string
	}{
		{"", nil, "unrecognized event id [], input: []"},

		{"START john", StartEvent{Name: "john"}, ""},
		{"START", nil, "1 argument(s) required, input: [START]"},

		{"WALK 5 1", WalkEvent{5, 1}, ""},
		{"WALK A 1", nil, "cannot read int argument [A], input: [WALK A 1]"},

		{"SHOOT 2 2", ShootEvent{2, 2}, ""},
		{"SHOOT 2 B", nil, "cannot read int argument [B], input: [SHOOT 2 B]"},
		{"MISS", MissEvent{}, ""},
		{"BOOM john night-king", BoomEvent{"john", "night-king"}, ""},
		{"VICTORY Petras", VictoryEvent{"Petras"}, ""},
		{"DEFEAT", DefeatEvent{}, ""},
	}

	for _, table := range tables {
		event, err := Unmarshal(table.msg)

		if table.err != "" {
			assert.EqualError(err, table.err)
		} else {
			assert.NoError(err)
		}

		assert.Equal(table.event, event)
	}
}
