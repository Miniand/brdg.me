package sushizock

import "testing"

const (
	Mick  = "Mick"
	Steve = "Steve"
)

func TestStart(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{Mick, Steve}); err != nil {
		t.Fatal(err)
	}
}
