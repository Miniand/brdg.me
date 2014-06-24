package for_sale

import "testing"

func TestStart(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{"mick", "steve", "bj"}); err != nil {
		t.Fatal(err)
	}
}
