package parser

import (
	"testing"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/texas_holdem"
	"github.com/Miniand/brdg.me/server/scommand"
)

func TestParseFrom(t *testing.T) {
	tests := map[string]string{
		"beefsack@gmail.com":                       "beefsack@gmail.com",
		"beefsack@gmail.com <Michael Alexander>":   "beefsack@gmail.com",
		" Michael Alexander <beefsack@gmail.com> ": "beefsack@gmail.com",
		"    beefsack@gmail.com blah":              "beefsack@gmail.com",
	}
	for source, target := range tests {
		from := ParseFrom(source)
		if from != target {
			t.Error("Could not detect email \"" + target + "\" from \"" +
				source + "\"")
		}
	}
}

func TestParseSubject(t *testing.T) {
	if ParseSubject("blah blah egg 52a7c891-e74d-463e-a47e-5c712a3dd439 art moo") !=
		"52a7c891-e74d-463e-a47e-5c712a3dd439" {
		t.Error("Unable to find objectid in subject")
	}
}

func TestParseBody(t *testing.T) {
	body := `

    pass  
  take    	 5
 
Kind regards,
Bob`
	ParseBody(body)
}

// @see https://github.com/Miniand/brdg.me/issues/22
func TestTexasHoldemRaiseBelowMin(t *testing.T) {
	g := &texas_holdem.Game{}
	err := g.Start([]string{"beefsack@gmail.com", "baconheist@gmail.com",
		"striker203@gmail.com"})
	if err != nil {
		t.Fatal(err)
	}
	commands := append(g.Commands(), scommand.Commands("")...)
	_, err = command.CallInCommands(g.WhoseTurn()[0], g, "raise 1", commands)
	if err == nil || err.Error() == "" {
		t.Fatal("Did not get an error!")
	}
}

func TestDecodeQuotedPrintable(t *testing.T) {
	expected := `If you believe that truth=beauty, then surely mathematics is the most beautiful branch of philosophy.`
	actual := DecodeQuotedPrintable("If you believe that truth=3Dbeauty, then surely =\r\nmathematics is the most beautiful branch of philosophy.")
	if actual != expected {
		t.Fatalf(`Expected: "%s" but got "%s"`, expected, actual)
	}
}

func TestValidateEmail(t *testing.T) {
	if !ValidateEmail("beefsack@gmail.com") {
		t.Error("beefsack@gmail.com should be valid")
	}
	if ValidateEmail("   fdsafad@egg.com") {
		t.Error("Shouldn't validate with whitespace")
	}
	if ValidateEmail("shonkydonk") {
		t.Error("Shouldn't validate ridiculous emails")
	}
}
