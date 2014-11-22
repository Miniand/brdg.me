package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardGiantBrain(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(names))
	g.Boards[Mick].Cards = []CardBase{&CardGiantBrain{}}
	g.RollPhase()
	cmd(t, g, Mick, "roll 1 2 3")
	cmd(t, g, Mick, "roll 1 2 3")
	cmd(t, g, Mick, "roll 1 2 3")
	cmd(t, g, Mick, "done")
}
