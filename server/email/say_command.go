package email

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
)

type SayCommand struct {
	gameId string
}

func (sc SayCommand) Parse(input string) []string {
	return command.ParseNamedCommand("say", input)
}

func (sc SayCommand) CanCall(player string, context interface{}) bool {
	g, ok := context.(game.Playable)
	if !ok || g.IsFinished() {
		return false
	}
	waitingOn := g.WhoseTurn()
	for _, p := range waitingOn {
		if p == player {
			// No point poking yourself!
			return false
		}
	}
	return true
}

func (sc SayCommand) Call(player string, context interface{},
	args []string) (string, error) {
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("No game was passed in")
	}
	g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(`%s says: %s`,
		render.PlayerNameInPlayers(player, g.PlayerList()),
		strings.Join(args[1:], " "))))
	if g.IsFinished() {
		// Just send it out to everyone.
		otherPlayers := []string{}
		for _, p := range g.PlayerList() {
			if p != player {
				otherPlayers = append(otherPlayers, p)
			}
		}
		CommunicateGameTo(sc.gameId, g, otherPlayers, "", false)
	}
	return "", nil
}

func (sc SayCommand) Usage(player string, context interface{}) string {
	return "{{b}}say ##{{_b}} to send a message to the other players, eg. {{b}}say hello!{{_b}}"
}
