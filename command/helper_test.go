package command

import (
	"testing"
)

func TestParseNamedCommandRangeArgs(t *testing.T) {
	args := ParseNamedCommandRangeArgs("bob", 2, 5, "bob egg cheese bacon")
	if args == nil {
		t.Fatal("Couldn't match command")
	}
	if len(args) < 2 {
		t.Fatal("Didn't get at least two args back", args)
	}
	actualArgs := ExtractNamedCommandArgs(args)
	if len(actualArgs) != 3 {
		t.Log(actualArgs)
		t.Fatal("Didn't get three actual args")
	}
	if actualArgs[0] != "egg" {
		t.Log(actualArgs[0])
		t.Fatal("actualArgs[0] wasn't 'egg', got:", actualArgs[0])
	}
	if actualArgs[1] != "cheese" {
		t.Log(actualArgs[1])
		t.Fatal("actualArgs[1] wasn't 'cheese', got:", actualArgs[1])
	}
	if actualArgs[2] != "bacon" {
		t.Log(actualArgs[2])
		t.Fatal("actualArgs[2] wasn't 'bacon', got:", actualArgs[2])
	}
}
