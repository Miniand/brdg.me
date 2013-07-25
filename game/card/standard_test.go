package card

import (
	"testing"
)

func TestRenderStandard52(t *testing.T) {
	c := SuitRankCard{
		Suit: STANDARD_52_SUIT_CLUBS,
		Rank: STANDARD_52_RANK_ACE,
	}
	expected := `{{c "black"}}♣A{{_c}}`
	output := c.RenderStandard52()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "black"}}♣A{{_c}} `
	output = c.RenderStandard52FixedWidth()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	c = SuitRankCard{
		Suit: STANDARD_52_SUIT_DIAMONDS,
		Rank: STANDARD_52_RANK_10,
	}
	expected = `{{c "red"}}♦10{{_c}}`
	output = c.RenderStandard52()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "red"}}♦10{{_c}}`
	output = c.RenderStandard52FixedWidth()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	c = SuitRankCard{
		Suit: STANDARD_52_SUIT_HEARTS,
		Rank: STANDARD_52_RANK_KING,
	}
	expected = `{{c "red"}}♥K{{_c}}`
	output = c.RenderStandard52()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "red"}}♥K{{_c}} `
	output = c.RenderStandard52FixedWidth()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	c = SuitRankCard{
		Suit: STANDARD_52_SUIT_SPADES,
		Rank: STANDARD_52_RANK_QUEEN,
	}
	expected = `{{c "black"}}♠Q{{_c}}`
	output = c.RenderStandard52()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "black"}}♠Q{{_c}} `
	output = c.RenderStandard52FixedWidth()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "gray"}}##{{_c}}`
	output = RenderStandard52Hidden()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
	expected = `{{c "gray"}}##{{_c}} `
	output = RenderStandard52HiddenFixedWidth()
	if output != expected {
		t.Error("Expected", expected, "but got", output)
	}
}

func TestAceHigh(t *testing.T) {
	d := Standard52DeckAceHigh()
	if d[0].(SuitRankCard).RankValue() <= STANDARD_52_RANK_KING {
		t.Fatal("Expected ace value to be higher than king")
	}
}
