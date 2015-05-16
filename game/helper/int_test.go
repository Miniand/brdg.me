package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntSort(t *testing.T) {
	unsorted := []int{5, 3, 4, 2}
	sorted := IntSort(unsorted)
	assert.Equal(t, []int{2, 3, 4, 5}, sorted)
	// Make sure the original slice wasn't modified.
	assert.Equal(t, []int{5, 3, 4, 2}, unsorted)
}

func TestIntReverse(t *testing.T) {
	ints := []int{5, 3, 4, 2}
	reversed := IntReverse(ints)
	assert.Equal(t, []int{2, 4, 3, 5}, reversed)
	// Make sure the original slice wasn't modified.
	assert.Equal(t, []int{5, 3, 4, 2}, ints)
}

func TestIntShuffle(t *testing.T) {
	unshuffled := []int{1, 2, 3, 4, 5}
	ok := false
	for i := 0; i < 1000; i++ {
		shuffled := IntShuffle(unshuffled)
		if !reflect.DeepEqual(unshuffled, shuffled) {
			ok = true
			break
		}
	}
	assert.True(t, ok)
	assert.Equal(t, unshuffled, []int{1, 2, 3, 4, 5})
}

func TestIntFind(t *testing.T) {
	k, ok := IntFind(4, []int{1, 2, 5, 6})
	assert.False(t, ok)
	k, ok = IntFind(3, []int{1, 2, 3, 3, 5})
	assert.True(t, ok)
	assert.Equal(t, 2, k)
}

func TestIntTally(t *testing.T) {
	assert.Equal(t, map[int]int{
		1: 4,
		2: 2,
		3: 1,
	}, IntTally([]int{1, 2, 3, 2, 1, 1, 1}))
}

func TestIntCount(t *testing.T) {
	assert.Equal(t, 3, IntCount(1, []int{1, 2, 3, 2, 1, 1}))
}

func TestIntFlatten(t *testing.T) {
	assert.Equal(t, []int{1, 1, 1, 2, 5, 5}, IntSort(IntFlatten(map[int]int{
		1: 3,
		2: 1,
		3: 0,
		4: -10,
		5: 2,
	})))
}
