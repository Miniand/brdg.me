package helper

import (
	"math/rand"
	"sort"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func IntSort(ints []int) []int {
	if ints == nil {
		return nil
	}
	sorted := make([]int, len(ints))
	copy(sorted, ints)
	sort.Ints(sorted)
	return sorted
}

func IntShuffle(ints []int) []int {
	if ints == nil {
		return nil
	}
	l := len(ints)
	shuffled := make([]int, l)
	for k, i := range rnd.Perm(l) {
		shuffled[k] = i
	}
	return shuffled
}

func IntFind(needle int, haystack []int) (index int, ok bool) {
	if haystack == nil {
		return
	}
	for k, i := range haystack {
		if i == needle {
			return k, true
		}
	}
	return
}

func IntTally(ints []int) map[int]int {
	if ints == nil {
		return nil
	}
	tally := map[int]int{}
	for _, i := range ints {
		tally[i]++
	}
	return tally
}

func IntCount(needle int, haystack []int) int {
	if haystack == nil {
		return 0
	}
	sum := 0
	for _, i := range haystack {
		if i == needle {
			sum++
		}
	}
	return sum
}

func IntFlatten(ints map[int]int) []int {
	if ints == nil {
		return nil
	}
	flat := []int{}
	for i, n := range ints {
		for j := 0; j < n; j++ {
			flat = append(flat, i)
		}
	}
	return flat
}
