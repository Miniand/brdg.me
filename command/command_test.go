package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type CommandTest struct{}

func (c CommandTest) Name() string {
	return "test"
}
func (c CommandTest) Call(
	player string,
	context interface{},
	input *Reader,
) (string, error) {
	return "tessssst", nil
}
func (c CommandTest) Usage(player string, context interface{}) string {
	return "Fart"
}

func TestCallInCommands(t *testing.T) {
	output, err := CallInCommands("bob", nil, `test
		  test`, []Command{CommandTest{}})
	assert.NoError(t, err)
	assert.Equal(t, "tessssst\ntessssst", output)
}
