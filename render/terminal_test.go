package render

import (
	"github.com/beefsack/brdg.me/game/tic_tac_toe"
	"testing"
)

func TestTerminalRender(t *testing.T) {
	g := &tic_tac_toe.Game{}
	output, err := RenderTerminal(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`, g)
	if err != nil {
		t.Error(err)
		return
	}
	if output != "\x1b[34m\x1b[34;1mhello\x1b[34m\x1b[0m" {
		t.Error("Output was", output)
		return
	}
}
