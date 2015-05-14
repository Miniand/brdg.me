package jaipur

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeck(t *testing.T) {
	assert.Len(t, Deck(), 55)
}
