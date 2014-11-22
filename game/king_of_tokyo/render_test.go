package king_of_tokyo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderAllCards(t *testing.T) {
	assert.NotEqual(t, "", RenderCardTable(Deck()))
}
