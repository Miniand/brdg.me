package alhambra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCard(t *testing.T) {
	c, err := ParseCard("R10")
	assert.NoError(t, err)
	assert.Equal(t, Card{CurrencyRed, 10}, c)
}
