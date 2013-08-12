package command

import (
	"testing"
)

func TestBasicCommand(t *testing.T) {
	c := BasicCommand{
		Name:    "play",
		NumArgs: 2,
		CanCallFunc: func(player string, context interface{}) bool {
			return true
		},
		CallFunc: func(player string, context interface{}, args []string) error {
			if len(args) != 2 {
				t.Fatal("There aren't two args")
			}
			if args[0] != "egg" {
				t.Fatal("First argument isn't egg")
			}
			if args[1] != "bacon" {
				t.Fatal("First argument isn't bacon")
			}
			return nil
		},
	}
	args := c.Parse("play egg bacon")
	if args == nil {
		t.Fatal("Args shouldn't be nil")
	}
	if args[1] != "egg" {
		t.Fatalf("Second argument wasn't egg, got: '%s'", args[1])
	}
	if args[2] != "bacon" {
		t.Fatalf("Third argument wasn't bacon, got: '%s'", args[2])
	}
	args = c.Parse("play egg")
	if args != nil {
		t.Fatal("Somehow successfully parsed command when not correct num args")
	}
}
