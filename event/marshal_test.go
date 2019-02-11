package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshal(t *testing.T) {
	assert := assert.New(t)
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
		msg, err := Marshal(table.event)
		assert.Equal(table.msg, msg, err)
	}
}
