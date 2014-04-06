package render

import (
	"fmt"
	"testing"
)

func TestReadPixel(t *testing.T) {
	cases := map[string][]string{
		"blah":                           []string{"b", "", "", "lah"},
		`{{c "red"}}YA{{_c}}`:            []string{"Y", `{{c "red"}}`, "", "A{{_c}}"},
		`{{c "red"}}{{b}}YA{{_b}}{{_c}}`: []string{"Y", `{{c "red"}}{{b}}`, "", "A{{_b}}{{_c}}"},
		`{{c "red"}}{{b}}Y{{_b}}{{_c}}`:  []string{"Y", `{{c "red"}}{{b}}`, `{{_b}}{{_c}}`, ""},
		`{{c "red"}}Y{{b}}A{{_b}}{{_c}}`: []string{"Y", `{{c "red"}}`, "", "{{b}}A{{_b}}{{_c}}"},
	}
	for input, expected := range cases {
		actualPixel, actualPrefix, actualSuffix, actualRemaining := readPixel(input)
		if expected[0] != actualPixel {
			t.Errorf("For %s, expected pixel to be %s but got %s",
				input, expected[0], actualPixel)
		}
		if expected[1] != actualPrefix {
			t.Errorf("For %s, expected prefix to be %s but got %s",
				input, expected[1], actualPrefix)
		}
		if expected[2] != actualSuffix {
			t.Errorf("For %s, expected suffix to be %s but got %s",
				input, expected[2], actualSuffix)
		}
		if expected[3] != actualRemaining {
			t.Errorf("For %s, expected remaining to be %s but got %s",
				input, expected[3], actualRemaining)
		}
	}
}

func ExampleCanvas_Render_positioning() {
	c := NewCanvas()
	c.Draw(-2, -4, "Hello...")
	c.Draw(1, -3, "World")
	fmt.Println(c.Render())
	// Output:
	// Hello...
	//    World
}

func ExampleCanvas_Render_tags() {
	c := NewCanvas()
	c.Draw(-1, -3, `{{b}}1{{_b}}`)
	c.Draw(1, -3, `{{c "green"}}2{{_c}}`)
	c.Draw(0, -2, `3`)
	fmt.Println(c.Render())
	// Output:
	// {{b}}1{{_b}} {{c "green"}}2{{_c}}
	//  3
}

func ExampleCanvas_Render_overwrite() {
	c := NewCanvas()
	c.Draw(-2, -4, `{{b}}Egg{{_b}}`)
	c.Draw(-3, -4, `Bacon`)
	fmt.Println(c.Render())
	// Output:
	// Bacon
}
