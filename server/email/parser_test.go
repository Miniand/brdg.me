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
	ParseBody(body)
}
