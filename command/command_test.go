package command

import (
	"github.com/beefsack/brdg.me/game"
	"regexp"
	"testing"
)

type CommandTest struct{}

func (c CommandTest) Parse(input string) []string {
	return regexp.MustCompile(`(?im)^\s*test\b\s$*`).FindStringSubmatch(input)
}
func (c CommandTest) CanCall(player string, g *game.Playable) bool {
	return true
}
func (c CommandTest) Call(player string, g *game.Playable, args []string) error {
	return nil
}
func (c CommandTest) Usage(player string, g *game.Playable) string {
	return "Fart"
}

func TestCommandInterface(t *testing.T) {
	c := &CommandTest{}
	args := c.Parse("Testicle")
	if args != nil {
		t.Fatal("args should be empty")
	}
	args = c.Parse(" test ")
	t.Logf("%#v", args)
	if args == nil {
		t.Fatalf("args shouldn't be empty")
	}
	if args[0] != " test " {
		t.Fatal("expected first arg to be test, got:", args[0])
	}
}

func TestCallInCommands(t *testing.T) {
	err := CallInCommands("bob", nil, `test
		  test`, []Command{CommandTest{}})
	if err != nil {
		t.Fatal(err)
	}
}
