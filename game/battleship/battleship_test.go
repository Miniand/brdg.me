package battleship

import (
	"fmt"
	"github.com/Miniand/brdg.me/command"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	MICK = iota
	STEVE
)

var names = map[int]string{
	MICK:  "Mick",
	STEVE: "Steve",
}

func mockGame(t *testing.T) *Game {
	players := []string{
		names[MICK],
		names[STEVE],
	}
	game := &Game{}
	err := game.Start(players)
	if err != nil {
		t.Fatal(err)
	}
	return game
}

func cloneGame(g *Game) *Game {
	newG := &Game{}
	data, err := g.Encode()
	if err != nil {
		panic(err.Error())
	}
	if err := newG.Decode(data); err != nil {
		panic(err.Error())
	}
	return newG
}

func TestNew(t *testing.T) {
	Convey("Given a new game, it should not fail", t, func() {
		mockGame(t)
	})
}

func TestParseShip(t *testing.T) {
	expectedMap := map[string]int{
		"sUB":     SHIP_SUBMARINE,
		" subma":  SHIP_SUBMARINE,
		"cru":     SHIP_CRUISER,
		"CAR":     SHIP_CARRIER,
		"destroy": SHIP_DESTROYER,
		"bat":     SHIP_BATTLESHIP,
	}
	for input, expected := range expectedMap {
		Convey(fmt.Sprintf("Given we try to parse %s", input), t, func() {
			input := input
			expected := expected
			actual, err := ParseShip(input)
			Convey("It should not error", func() {
				So(err, ShouldBeNil)
			})
			Convey(fmt.Sprintf("It should parse as a %s", shipNames[expected]), func() {
				So(actual, ShouldEqual, expected)
			})
		})
	}
}

func TestParseDirection(t *testing.T) {
	expectedMap := map[string]int{
		"up":    DIRECTION_UP,
		"DOWN":  DIRECTION_DOWN,
		"LEft":  DIRECTION_LEFT,
		"right": DIRECTION_RIGHT,
		"u":     DIRECTION_UP,
		"d":     DIRECTION_DOWN,
		"l":     DIRECTION_LEFT,
		"r":     DIRECTION_RIGHT,
		"north": DIRECTION_UP,
		"SOUTH": DIRECTION_DOWN,
		"west":  DIRECTION_LEFT,
		"easT":  DIRECTION_RIGHT,
		"n":     DIRECTION_UP,
		"S":     DIRECTION_DOWN,
		"w":     DIRECTION_LEFT,
		"E":     DIRECTION_RIGHT,
	}
	for input, expected := range expectedMap {
		Convey(fmt.Sprintf("Given we try to parse %s", input), t, func() {
			input := input
			expected := expected
			actual, err := ParseDirection(input)
			Convey("It should not error", func() {
				So(err, ShouldBeNil)
			})
			Convey(fmt.Sprintf("It should parse as %s", directionNames[expected]), func() {
				So(actual, ShouldEqual, expected)
			})
		})
	}
}

func TestLocationName(t *testing.T) {
	input := [][2]int{
		[2]int{Y_B, X_3},
		[2]int{Y_A, X_1},
		[2]int{Y_J, X_10},
	}
	output := []string{
		"B3",
		"A1",
		"J10",
	}
	for index, input := range input {
		Convey(fmt.Sprintf("Given we try to output %d and %d", input[0], input[1]), t, func() {
			expected := output[index]
			actual := LocationName(input[0], input[1])
			Convey(fmt.Sprintf("It should output %s", expected), func() {
				So(actual, ShouldEqual, expected)
			})
		})
	}
}

func TestParseLocation(t *testing.T) {
	expectedMap := map[string][]int{
		"b3": []int{Y_B, X_3},
		"C7": []int{Y_C, X_7},
	}
	for input, expected := range expectedMap {
		Convey(fmt.Sprintf("Given we try to parse %s", input), t, func() {
			input := input
			expected := expected
			actualY, actualX, err := ParseLocation(input)
			Convey("It should not error", func() {
				So(err, ShouldBeNil)
			})
			Convey(fmt.Sprintf("It should parse as %s", LocationName(
				expected[0], expected[1])), func() {
				So(actualY, ShouldEqual, expected[0])
				So(actualX, ShouldEqual, expected[1])
			})
		})
	}
}

func TestGame(t *testing.T) {
	g := mockGame(t)
	// Both players place
	if len(g.WhoseTurn()) != 2 {
		t.Fatal("Both players should be placing")
	}
	_, err := command.CallInCommands(names[MICK], g,
		"place sub b3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[MICK], g,
		"place car c3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[MICK], g,
		"place des d3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[MICK], g,
		"place cru e3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[MICK], g,
		"place bat f3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[STEVE], g,
		"place sub b3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[STEVE], g,
		"place car c3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[STEVE], g,
		"place des d3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[STEVE], g,
		"place cru e3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[STEVE], g,
		"place bat f3 right", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
	_, err = command.CallInCommands(names[MICK], g,
		"shoot b3", g.Commands())
	if err != nil {
		t.Fatal(err)
	}
}
