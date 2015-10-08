package seven_wonders_duel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCards(t *testing.T) {
	for cId, c := range Cards {
		assert.NotEmpty(t, c.Name, fmt.Sprintf("card %d is missing a name", cId))
		assert.NotEmpty(t, c.Id, fmt.Sprintf("card %s is missing an Id", c.Name))
		assert.Equal(t, cId, c.Id, fmt.Sprintf("card %s doesn't have an Id matching the card key (%d != %d)", c.Name, cId, c.Id))
		assert.NotEmpty(t, c.Type, fmt.Sprintf("card %s is missing a Type", c.Name))
	}
}
