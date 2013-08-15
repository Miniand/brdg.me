package command

import (
	"regexp"
	"testing"
)

type CommandTest struct{}

func (c CommandTest) Parse(input string) []string {
	return regexp.MustCompile(`(?im)^\s*test\b\s$*`).FindStringSubmatch(input)
}
func (c CommandTest) CanCall(player string, context interface{}) bool {
	return true
}
func (c CommandTest) Call(player string, context interface{},
	args []string) (string, error) {
	return "tessssst", nil
}
func (c CommandTest) Usage(player string, context interface{}) string {
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
	output, err := CallInCommands("bob", nil, `test
		  test`, []Command{CommandTest{}})
	if err != nil {
		t.Fatal(err)
	}
	if output != "tessssst" {
		t.Fatal("Expected output to be tessssst, got:", output)
	}
}
