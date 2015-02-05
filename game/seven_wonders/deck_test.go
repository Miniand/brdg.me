package seven_wonders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeckAge1(t *testing.T) {
	for p := 3; p <= 7; p++ {
		assert.Len(t, DeckAge1(p), p*7)
	}
}

func TestDeckAge2(t *testing.T) {
	for p := 3; p <= 7; p++ {
		assert.Len(t, DeckAge2(p), p*7)
	}
}

func TestDeckAge3(t *testing.T) {
	for p := 3; p <= 7; p++ {
		assert.Len(t, DeckAge3(p), p*7)
	}
}
