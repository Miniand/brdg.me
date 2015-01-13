package splendor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNobleCards(t *testing.T) {
	assert.Len(t, NobleCards(), 10)
}
