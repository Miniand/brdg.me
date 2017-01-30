package age_of_war

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Miniand/brdg.me/game/helper"
	"github.com/stretchr/testify/assert"
)

func TestGame_Start(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(helper.Players[:3]))
}

func TestGame_IsFinished(t *testing.T) {
	g := &Game{}
	f, err := os.Open("testdata/1e3f5cd4-ee30-4bc7-8f76-36c3fdc55476")
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.NoError(t, g.Decode(data))
	assert.Equal(t, []string{"baconheist@gmail.com"}, g.Winners())
}
