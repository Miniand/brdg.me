package starship_catan

import (
	"reflect"
	"testing"

	"github.com/Miniand/brdg.me/command"
)

const (
	Mick  = "Mick"
	Steve = "Steve"
)

func TestStart(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{Mick, Steve}); err != nil {
		t.Fatal(err)
	}
}

func TestChooseModule(t *testing.T) {
	g := &Game{}
	if err := g.Start([]string{Mick, Steve}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual([]string{Mick, Steve}, g.WhoseTurn()) {
		t.Fatal("Expected it to be both Mick and Steve's turn.")
	}
	if _, err := command.CallInCommands(Steve, g,
		"choose lo", g.Commands(Steve)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(map[int]int{
		ModuleLogistics: 1,
	}, g.PlayerBoards[1].Modules) {
		t.Fatal("Expected Steve to have a level 1 logistics module.")
	}
	if !reflect.DeepEqual([]string{Mick}, g.WhoseTurn()) {
		t.Fatal("Expected it to be only Mick's turn.")
	}
	if _, err := command.CallInCommands(Mick, g,
		"choose se", g.Commands(Mick)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(map[int]int{
		ModuleSensor: 1,
	}, g.PlayerBoards[0].Modules) {
		t.Fatal("Expected Mick to have a level 1 sensor module.")
	}
	if !reflect.DeepEqual([]string{Mick}, g.WhoseTurn()) {
		t.Fatal("Expected it to be only Mick's turn.")
	}
	if g.Phase == PhaseChooseModule {
		t.Fatal("It should no longer be the choose module phase.")
	}
}
