package hive

import "testing"

func TestStartGame(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{"Mick", "Steve"}); err != nil {
		t.Fatal(err)
	}
}
