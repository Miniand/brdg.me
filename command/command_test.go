package command

import (
	"regexp"
	"testing"
)

type CommandTest struct{}

func (c CommandTest) Name() string {
	return "test"
}
func (c CommandTest) Parse(input string) []string {
	return regexp.MustCompile(`(?im)^\s*test\b\s$*`).FindStringSubmatch(input)
}
func (c CommandTest) Call(
	player string,
	context interface{},
	input *Parser,
) (string, error) {
	return "tessssst", nil
}
func (c CommandTest) Usage(player string, context interface{}) string {
	return "Fart"
}

func TestCallInCommands(t *testing.T) {
	output, err := CallInCommands("bob", nil, `test
		  test`, []Command{CommandTest{}})
	if err != nil {
		t.Fatal(err)
	}
	if output != "tessssst" {
		t.Fatal("Expected output to be tessssst, got:", output)
	}
}
