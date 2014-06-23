package for_sale

import "testing"

func TestStart(t *testing.T) {
	g := &Game{}
	g.Start([]string{"mick", "steve", "bj"})
}
