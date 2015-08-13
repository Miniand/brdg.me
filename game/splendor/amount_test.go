package splendor

import (
	"testing"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/stretchr/testify/assert"
)

func TestCanAfford(t *testing.T) {
	assert.True(t, CanAfford(cost.Cost{
		Emerald: 2,
		Gold:    1,
	}, cost.Cost{
		Emerald: 3,
	}))
	assert.False(t, CanAfford(cost.Cost{
		Emerald: 2,
		Gold:    1,
	}, cost.Cost{
		Emerald: 4,
	}))
}
