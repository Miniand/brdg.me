package starship_catan

import "testing"

func TestStart(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{"mick", "steve"}); err != nil {
		t.Fatal(err)
	}
}
