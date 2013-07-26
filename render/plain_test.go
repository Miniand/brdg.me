package render

import (
	"github.com/beefsack/brdg.me/game/tic_tac_toe"
	"testing"
)

func TestPlainRender(t *testing.T) {
	g := &tic_tac_toe.Game{}
	output, err := RenderPlain(`{{c "blue"}}{{b}}hello{{_b}}{{_c}}`, g)
	if err != nil {
		t.Error(err)
		return
	}
	if output != `hello` {
		t.Error("Output was", output)
		return
	}
}
