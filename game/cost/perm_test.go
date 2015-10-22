package cost

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type permCase struct {
	Cost            Cost
	With            [][]Cost
	ExpectedCan     bool
	ExpectedCanWith [][]Cost
}

func deepTrim(costs []Cost) []Cost {
	trimmed := make([]Cost, len(costs))
	for i, c := range costs {
		trimmed[i] = c.Trim()
	}
	return trimmed
}
func doubleDeepTrim(costs [][]Cost) [][]Cost {
	trimmed := make([][]Cost, len(costs))
	for i, c := range costs {
		trimmed[i] = deepTrim(c)
	}
	return trimmed
}

func TestCanAffordPerm(t *testing.T) {
	for testName, testCase := range map[string]permCase{
		"Can afford with nothing": {
			Cost:            Cost{},
			With:            [][]Cost{},
			ExpectedCan:     true,
			ExpectedCanWith: [][]Cost{},
		},
		"Can't afford with nothing": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With:            [][]Cost{},
			ExpectedCan:     false,
			ExpectedCanWith: [][]Cost{},
		},
		"Can afford with single": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 6,
					},
				},
			},
			ExpectedCan: true,
			ExpectedCanWith: [][]Cost{
				{
					{
						TestRes1: 3,
						TestRes2: 4,
					},
				},
			},
		},
		"Can't afford with single": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 2,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan:     false,
			ExpectedCanWith: [][]Cost{},
		},
		"Can afford with multiple": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
				},
				{
					{
						TestRes1: 1,
						TestRes2: 3,
					},
				},
				{
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan: true,
			ExpectedCanWith: [][]Cost{
				{
					{
						TestRes1: 3,
						TestRes2: 0,
					},
					{
						TestRes1: 0,
						TestRes2: 3,
					},
					{
						TestRes1: 0,
						TestRes2: 1,
					},
				},
			},
		},
		"Can't afford with multiple": {
			Cost: Cost{
				TestRes1: 8,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
				},
				{
					{
						TestRes1: 1,
						TestRes2: 3,
					},
				},
				{
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan:     false,
			ExpectedCanWith: [][]Cost{},
		},
		"Can afford with perm": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
					{
						TestRes1: 5,
						TestRes2: 6,
					},
				},
			},
			ExpectedCan: true,
			ExpectedCanWith: [][]Cost{
				{
					{
						TestRes1: 3,
						TestRes2: 4,
					},
				},
			},
		},
		"Can't afford with perm": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
					{
						TestRes1: 5,
						TestRes2: 2,
					},
				},
			},
			ExpectedCan:     false,
			ExpectedCanWith: [][]Cost{},
		},
		"Can afford with multiple perm": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 4,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan: true,
			ExpectedCanWith: [][]Cost{
				{
					{
						TestRes1: 3,
						TestRes2: 0,
					},
					{
						TestRes1: 0,
						TestRes2: 4,
					},
				},
				{
					{
						TestRes1: 1,
						TestRes2: 4,
					},
					{
						TestRes1: 2,
						TestRes2: 0,
					},
				},
			},
		},
		"Can't afford with multiple perm": {
			Cost: Cost{
				TestRes1: 6,
				TestRes2: 7,
			},
			With: [][]Cost{
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
				{
					{
						TestRes1: 5,
						TestRes2: 0,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan:     false,
			ExpectedCanWith: [][]Cost{},
		},
		"Can afford with multiple perm some irrelevant": {
			Cost: Cost{
				TestRes1: 3,
				TestRes2: 1,
			},
			With: [][]Cost{
				{
					{
						TestRes3: 5,
					},
					{
						TestRes4: 1,
					},
				},
				{
					{
						TestRes1: 5,
						TestRes2: 2,
					},
					{
						TestRes1: 1,
						TestRes2: 5,
					},
				},
			},
			ExpectedCan: true,
			ExpectedCanWith: [][]Cost{
				{
					{
						TestRes1: 0,
						TestRes2: 0,
					},
					{
						TestRes1: 3,
						TestRes2: 1,
					},
				},
			},
		},
	} {
		actualCan, actualCanWith := CanAffordPerm(testCase.Cost, testCase.With)
		assert.Equal(t, testCase.ExpectedCan, actualCan,
			fmt.Sprintf("can incorrect for test case '%s'", testName))
		assert.Equal(
			t,
			doubleDeepTrim(testCase.ExpectedCanWith),
			doubleDeepTrim(actualCanWith),
			fmt.Sprintf("canWith incorrect for test case '%s'", testName),
		)
	}
}

func TestPerm(t *testing.T) {
	assert.Equal(t, []Cost{}, Perm([][]Cost{}))
	assert.Equal(t, []Cost{
		{1: 3, 2: 4},
	}, Perm([][]Cost{
		{
			{1: 3, 2: 1},
		},
		{
			{2: 3},
		},
	}))
	assert.Equal(t, []Cost{
		{1: 3, 2: 4},
		{1: 2, 2: 3, 3: 5},
	}, Perm([][]Cost{
		{
			{1: 3, 2: 1},
			{1: 2, 3: 5},
		},
		{
			{2: 3},
		},
	}))
	assert.Equal(t, []Cost{
		{1: 3, 2: 4},
		{1: 2, 2: 3, 3: 5},
		{1: 3, 2: 1, 3: 4},
		{1: 2, 3: 9},
	}, Perm([][]Cost{
		{
			{1: 3, 2: 1},
			{1: 2, 3: 5},
		},
		{
			{2: 3},
			{3: 4},
		},
	}))
}
