package acquire

import (
	"testing"
)

func TestStart(t *testing.T) {
	g := Game{}
	if err := g.Start([]string{"Mick", "Steve"}); err != nil {
		t.Fatal(err)
	}
}
