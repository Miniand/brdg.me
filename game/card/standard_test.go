package card

import (
	"testing"
)

func TestRenderStandard52(t *testing.T) {
	c := SuitValueCard{
		Suit:  STANDARD_52_SUIT_CLUBS,
		Value: STANDARD_52_VALUE_ACE,
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
	c = SuitValueCard{
		Suit:  STANDARD_52_SUIT_DIAMONDS,
		Value: STANDARD_52_VALUE_10,
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
	c = SuitValueCard{
		Suit:  STANDARD_52_SUIT_HEARTS,
		Value: STANDARD_52_VALUE_KING,
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
	c = SuitValueCard{
		Suit:  STANDARD_52_SUIT_SPADES,
		Value: STANDARD_52_VALUE_QUEEN,
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
