package lost_cities

import (
	"testing"
)

func TestStart(t *testing.T) {
	players := []string{"Mick", "Steve"}
	game := &Game{}
	err := game.Start(players)
	//err := errors.New("dunno wtf I'm doing")
	if err != nil {
		t.Error(err)
	}
}
