package starship_catan

import "testing"

func TestParseModule(t *testing.T) {
	_, err := ParseModule("FARTS!")
	if err == nil {
		t.Fatal("Expected not to find farts module.")
	}
	m, err := ParseModule("L")
	if err != nil {
		t.Fatal(err)
	}
	if m != ModuleLogistics {
		t.Fatal("Expected module to be logistics")
	}
	_, err = ParseModule("s")
	if err == nil {
		t.Fatal("Expected not to find unique module starting with s.")
	}
	m, err = ParseModule("se")
	if err != nil {
		t.Fatal(err)
	}
	if m != ModuleSensor {
		t.Fatal("Expected module to be sensor")
	}
}
