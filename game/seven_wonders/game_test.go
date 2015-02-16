package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestPlayerScienceVP(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Cards[Mick] = card.Deck{
		Cards[WonderStageBabylonA2],
	}
	assert.Equal(t, 1, g.PlayerScienceVP(Mick))
}

func TestScienceVPPerm(t *testing.T) {
	assert.Equal(t, [][]int{}, ScienceVPPerm([][]int{}))

	assert.Equal(t, [][]int{
		{FieldEngineering},
	}, ScienceVPPerm([][]int{
		{FieldEngineering},
	}))

	assert.Equal(t, [][]int{
		{FieldEngineering, FieldMathematics, FieldMathematics},
		{FieldEngineering, FieldTheology, FieldMathematics},
	}, ScienceVPPerm([][]int{
		{FieldEngineering},
		{FieldMathematics, FieldTheology},
		{FieldMathematics},
	}))
}

func TestScienceVP(t *testing.T) {
	assert.Equal(t, 0, ScienceVP([]int{}))

	assert.Equal(t, 10, ScienceVP([]int{
		FieldEngineering,
		FieldTheology,
		FieldMathematics,
	}))

	assert.Equal(t, 13, ScienceVP([]int{
		FieldEngineering,
		FieldEngineering,
		FieldTheology,
		FieldMathematics,
	}))

	assert.Equal(t, 16, ScienceVP([]int{
		FieldEngineering,
		FieldEngineering,
		FieldEngineering,
		FieldEngineering,
	}))

	assert.Equal(t, 26, ScienceVP([]int{
		FieldEngineering,
		FieldEngineering,
		FieldTheology,
		FieldTheology,
		FieldMathematics,
		FieldMathematics,
	}))
}
