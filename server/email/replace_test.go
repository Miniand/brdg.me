package main

import (
	"testing"
)

func TestReplaceCharactersWithCidImg(t *testing.T) {
	str := "This is a Dice!!!  ⚀!!!!"
	expected := `This is a Dice!!!  <img src="cid:charreplace9856@brdg.me" style="vertical-align:bottom;" alt="⚀" />!!!!`
	actual, found := ReplaceCharactersWithCidImg(str, CharacterReplacements)
	if actual != expected {
		t.Fatal("Expected", expected, "but got", actual)
	}
	if len(found) != 1 {
		t.Fatal("Expected found would be length 1, got", len(found))
	}
	if found[0] != '⚀' {
		t.Fatal("Expected found to be ⚀, got", found[0])
	}
}

func TestRuneCid(t *testing.T) {
	expected := "charreplace9856@brdg.me"
	actual := RuneCid('⚀')
	if actual != expected {
		t.Fatal("Expected", expected, "but got", actual)
	}
}

func TestRuneReplacementCidImg(t *testing.T) {
	expected := `<img src="cid:charreplace9856@brdg.me" style="vertical-align:bottom;" alt="⚀" />`
	actual := RuneReplacementCidImg('⚀')
	if actual != expected {
		t.Fatal("Expected", expected, "but got", actual)
	}
}
