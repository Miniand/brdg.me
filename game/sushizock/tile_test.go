package sushizock

import (
	"reflect"
	"testing"
)

func TestTiles_Remove(t *testing.T) {
	tiles := Tiles{
		{TileTypeBlue, 1},
		{TileTypeBlue, 2},
		{TileTypeBlue, 3},
		{TileTypeBlue, 4},
		{TileTypeBlue, 5},
	}
	actual, remaining := tiles.Remove(2)
	expected := Tile{TileTypeBlue, 3}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v to equal %#v", expected, actual)
	}
	expectedTiles := Tiles{
		{TileTypeBlue, 1},
		{TileTypeBlue, 2},
		{TileTypeBlue, 4},
		{TileTypeBlue, 5},
	}
	if !reflect.DeepEqual(expectedTiles, remaining) {
		t.Errorf("Expected %#v to equal %#v", expectedTiles, remaining)
	}
}
