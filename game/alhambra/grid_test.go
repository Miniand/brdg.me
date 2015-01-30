package alhambra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrid_IsValid_valid(t *testing.T) {
	g := parseGrid(t, `

 F

`)
	valid, reason := g.IsValid()
	assert.True(t, valid)
	assert.Empty(t, reason)
}

func TestGrid_IsValid_invalidNoFountain(t *testing.T) {
	g := parseGrid(t, `

 A|S

`)
	valid, reason := g.IsValid()
	assert.False(t, valid)
	assert.Equal(t, GridInvalidNoFountain, reason)
}

func TestGrid_IsValid_invalidWall(t *testing.T) {
	g := parseGrid(t, `

 A|S F

 A S

`)
	g[Vect{0, 0}].Walls[DirRight] = false // Remove the right wall from TL
	valid, reason := g.IsValid()
	assert.False(t, valid)
	assert.Equal(t, GridInvalidWall, reason)
}

func TestGrid_IsValid_invalidCannotWalk(t *testing.T) {
	g := parseGrid(t, `

 A A|S F
 - -
   S S

`)
	valid, reason := g.IsValid()
	assert.False(t, valid)
	assert.Equal(t, GridInvalidCannotWalk, reason)
}

func TestGrid_IsValid_invalidGap(t *testing.T) {
	g := parseGrid(t, `

 A S F
 
 A   A

   A A

`)
	valid, reason := g.IsValid()
	assert.False(t, valid)
	assert.Equal(t, GridInvalidGap, reason)
}

func TestGrid_LongestExtWall(t *testing.T) {
	assert.Equal(t, 5, parseGrid(t, `
+-
|A A A
     -+
 A A A|
     -+
 A A A|
 -----+-+
       A|
       -+
`).LongestExtWall())
	assert.Equal(t, 12, parseGrid(t, `
+-----+
|A A A|
|    -+
|A A A|
|  ---+
|A A A|
+-----+
`).LongestExtWall())
}
