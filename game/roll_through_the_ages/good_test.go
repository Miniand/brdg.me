package roll_through_the_ages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoodMaximum(t *testing.T) {
	assert.Equal(t, 8, GoodMaximum(GoodWood))
	assert.Equal(t, 7, GoodMaximum(GoodStone))
	assert.Equal(t, 6, GoodMaximum(GoodPottery))
	assert.Equal(t, 5, GoodMaximum(GoodCloth))
	assert.Equal(t, 4, GoodMaximum(GoodSpearhead))
}

func TestGoodValue(t *testing.T) {
	assert.Equal(t, 1, GoodValue(GoodWood, 1))
	assert.Equal(t, 10, GoodValue(GoodWood, 4))
	assert.Equal(t, 36, GoodValue(GoodWood, 8))

	assert.Equal(t, 2, GoodValue(GoodStone, 1))
	assert.Equal(t, 12, GoodValue(GoodStone, 3))
	assert.Equal(t, 56, GoodValue(GoodStone, 7))

	assert.Equal(t, 3, GoodValue(GoodPottery, 1))
	assert.Equal(t, 18, GoodValue(GoodPottery, 3))
	assert.Equal(t, 63, GoodValue(GoodPottery, 6))

	assert.Equal(t, 4, GoodValue(GoodCloth, 1))
	assert.Equal(t, 24, GoodValue(GoodCloth, 3))
	assert.Equal(t, 60, GoodValue(GoodCloth, 5))

	assert.Equal(t, 5, GoodValue(GoodSpearhead, 1))
	assert.Equal(t, 30, GoodValue(GoodSpearhead, 3))
	assert.Equal(t, 50, GoodValue(GoodSpearhead, 4))
}
