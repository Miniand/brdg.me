package cost

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestRes1 = iota
	TestRes2
	TestRes3
	TestRes4
	TestRes5
)

func TestCost_Clone(t *testing.T) {
	c1 := Cost{
		TestRes1: 4,
		TestRes2: 5,
	}
	c2 := c1.Clone()
	c2[TestRes1] = 10
	assert.Equal(t, Cost{
		TestRes1: 4,
		TestRes2: 5,
	}, c1)
}

func TestCost_Add(t *testing.T) {
	c1 := Cost{
		TestRes1: 4,
		TestRes2: 5,
	}
	c2 := Cost{
		TestRes1: 3,
	}
	assert.Equal(t, Cost{
		TestRes1: 7,
		TestRes2: 5,
	}, c1.Add(c2))
}

func TestCost_Inv(t *testing.T) {
	c1 := Cost{
		TestRes1: -3,
		TestRes2: 6,
	}
	assert.Equal(t, Cost{
		TestRes1: 3,
		TestRes2: -6,
	}, c1.Inv())
}

func TestCost_Sub(t *testing.T) {
	c1 := Cost{
		TestRes1: 2,
		TestRes2: 3,
	}
	c2 := Cost{
		TestRes1: 1,
		TestRes2: 4,
	}
	assert.Equal(t, Cost{
		TestRes1: 1,
		TestRes2: -1,
	}, c1.Sub(c2))
}

func TestCost_SignSplit(t *testing.T) {
	c := Cost{
		TestRes1: 4,
		TestRes2: -5,
	}
	pos, neg := c.PosNeg()
	assert.Equal(t, Cost{
		TestRes1: 4,
	}, pos)
	assert.Equal(t, Cost{
		TestRes2: -5,
	}, neg)
}

func TestCost_CanAfford(t *testing.T) {
	c := Cost{
		TestRes1: 3,
		TestRes2: 4,
	}
	assert.True(t, c.CanAfford(Cost{
		TestRes1: 2,
		TestRes2: 4,
	}))
	assert.False(t, c.CanAfford(Cost{
		TestRes1: 5,
		TestRes2: 4,
	}))
}

func TestCost_Take(t *testing.T) {
	c := Cost{
		TestRes1: 4,
		TestRes2: 5,
	}
	assert.Equal(t, Cost{
		TestRes1: 4,
	}, c.Take(TestRes1))
}

func TestCost_Drop(t *testing.T) {
	c := Cost{
		TestRes1: 4,
		TestRes2: 5,
	}
	assert.Equal(t, Cost{
		TestRes2: 5,
	}, c.Drop(TestRes1))
}

func TestCost_Ints(t *testing.T) {
	c := Cost{
		TestRes1: 2,
		TestRes2: 1,
		TestRes3: 3,
	}
	assert.Equal(t, []int{
		TestRes1,
		TestRes1,
		TestRes2,
		TestRes3,
		TestRes3,
		TestRes3,
	}, c.Ints())
}

func TestCost_Sum(t *testing.T) {
	c := Cost{
		TestRes1: 2,
		TestRes2: 1,
		TestRes3: 3,
	}
	assert.Equal(t, 6, c.Sum())
}
