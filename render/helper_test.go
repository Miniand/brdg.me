package render

import (
	"github.com/beefsack/brdg.me/game/tic_tac_toe"
	"testing"
)

func TestPlayerColour(t *testing.T) {
	if PlayerColour(0) != "green" {
		t.Fatal("Expected first player to be green")
	}
	if PlayerColour(9) != "red" {
		t.Fatal("Expected tenth player to be red")
	}
}

func TestPlayerName(t *testing.T) {
	if PlayerName(1, "bob") != `{{b}}{{c "red"}}bob{{_c}}{{_b}}` {
		t.Fatal("bob didn't render bold and red")
	}
}

func TestPadded(t *testing.T) {
	g := &tic_tac_toe.Game{}
	text, err := Padded("{{b}}你好{{_b}}", 5, g)
	if err != nil {
		t.Fatal(err)
	}
	if text != "{{b}}你好{{_b}}   " {
		t.Fatal("Expected 你好 to gain three spaces, got:", text)
	}
}
