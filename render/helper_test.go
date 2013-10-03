package render

import (
	"testing"
)

func TestPlayerColour(t *testing.T) {
	if PlayerColour(0) != "green" {
		t.Fatal("Expected first player to be green")
	}
	if PlayerColour(9) != "red" {
		t.Fatal("Expected tenth player to be red")
	}
}

func TestPlayerName(t *testing.T) {
	if PlayerName(1, "bob") != `{{b}}{{c "red"}}bob{{_c}}{{_b}}` {
		t.Fatal("bob didn't render bold and red")
	}
}

func TestPadded(t *testing.T) {
	text, err := Padded("{{b}}你好{{_b}}", 5)
	if err != nil {
		t.Fatal(err)
	}
	if text != "{{b}}你好{{_b}}   " {
		t.Fatal("Expected 你好 to gain three spaces, got:", text)
	}
}

func TestTable(t *testing.T) {
	output, err := Table([][]string{}, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if output != "" {
		t.Fatal("Output wasn't blank, got:", output)
	}
	output, err = Table([][]string{
		[]string{"{{b}}Five{{_b}}", "One"},
		[]string{"Twenty"},
	}, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if output != `{{b}}Five{{_b}}   One
Twenty` {
		t.Fatal("Output wasn't correct, got:", output)
	}
}

func TestCommaList(t *testing.T) {
	var expected, actual string
	expected = "cheese"
	actual = CommaList([]string{"cheese"})
	if actual != expected {
		t.Fatal("Expected", expected, "got", actual)
	}
	expected = "cheese and burger"
	actual = CommaList([]string{"cheese", "burger"})
	if actual != expected {
		t.Fatal("Expected", expected, "got", actual)
	}
	expected = "fart, cheese and burger"
	actual = CommaList([]string{"fart", "cheese", "burger"})
	if actual != expected {
		t.Fatal("Expected", expected, "got", actual)
	}
	expected = "bogs, fart, cheese and burger"
	actual = CommaList([]string{"bogs", "fart", "cheese", "burger"})
	if actual != expected {
		t.Fatal("Expected", expected, "got", actual)
	}
}
