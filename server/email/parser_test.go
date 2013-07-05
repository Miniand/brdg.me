package main

import (
	"testing"
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
	if ParseSubject("blah blah egg 2847fbac9506adcd4587edbc art moo") !=
		"2847fbac9506adcd4587edbc" {
		t.Error("Unable to find objectid in subject")
	}
}

func TestParseBody(t *testing.T) {
	body := `

    pass  
  take    	 5
 
Kind regards,
Bob`
	commands := ParseBody(body)
	if len(commands) != 2 {
		t.Error("Command count incorrect, expected 2, got", len(commands))
	}
	if len(commands[0]) != 1 {
		t.Error("Expected one part for the first command, got",
			len(commands[0]))
	}
	if commands[0][0] != "pass" {
		t.Error("Expected first command to be pass, got", commands[0][0])
	}
}
