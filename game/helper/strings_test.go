package helper

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchStringInStrings(t *testing.T) {
	for _, test := range []struct {
		input    string
		strs     []string
		expected int
	}{
		{"blah", []string{"BLAHE", "BLAH"}, 1},
		{"blah", []string{"BLAHE", "BLAHr"}, -1},
		{"M", []string{"mick", "steve"}, 0},
		{"M", []string{"mick", "steve", "mork"}, -1},
		{"Mi", []string{"mick", "steve", "mork"}, 0},
	} {
		m, err := MatchStringInStrings(test.input, test.strs)
		if err != nil && test.expected != -1 {
			t.Errorf(`Expected %s to match %s in (%s) but got error: %s`,
				test.input, test.expected, strings.Join(test.strs, ", "), err)
		}
		if err == nil && m != test.expected {
			t.Errorf(`Expected %s to match %s in (%s)`,
				test.input, test.expected, strings.Join(test.strs, ", "))
		}
	}
}

func TestNumberStr(t *testing.T) {
	assert.Equal(t, "fourty three billion two hundred and eighty one million five hundred and twenty three thousand six hundred and eighty one", NumberStr(43281523681))
}
