package scommand

import (
	"errors"
	"fmt"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/helper"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
)

type ConcedeCommand struct {
	gameModel *model.GameModel
}

func (c ConcedeCommand) Parse(input string) []string {
	return command.ParseNamedCommandRangeArgs("concede", 0, -1, input)
}

func CanInitiateConcedeVote(player string, gm *model.GameModel) bool {
	if gm.IsFinished || gm.ConcedeVote != nil {
		return false
	}
	if gm.EliminatedPlayerList != nil {
		// Eliminated players can't vote to concede
		for _, p := range gm.EliminatedPlayerList {
			if p == player {
				return false
			}
		}
	}
	return true
}

func CanConcedeVote(player string, gm *model.GameModel) bool {
	for _, p := range gm.RemainingConcedeVotePlayers() {
		if p == player {
			return true
		}
	}
	return false
}

func PassConcedeVote(gm *model.GameModel, g game.Playable) {
	gm.IsFinished = true
	gm.Winners = gm.ConcedePlayers
	gm.ConcedePlayers = nil
	gm.ConcedeVote = nil
	g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
		"{{b}}The game has been conceeded to %s{{_b}}",
		render.CommaList(render.PlayerNamesInPlayers(gm.Winners, gm.PlayerList)),
	)))
}

func FailConcedeVote(gm *model.GameModel, g game.Playable) {
	gm.ConcedePlayers = nil
	gm.ConcedeVote = nil
	g.GameLog().Add(log.NewPublicMessage("{{b}}The vote failed{{_b}}"))
}

func (c ConcedeCommand) CanCall(player string, context interface{}) bool {
	return CanInitiateConcedeVote(player, c.gameModel)
}

func (c ConcedeCommand) Call(player string, context interface{},
	args []string) (string, error) {
	if !c.CanCall(player, context) {
		return "", errors.New("you can't concede at the moment")
	}
	g, ok := context.(game.Playable)
	if !ok {
		return "", errors.New("no game was passed in")
	}

	playerList := c.gameModel.PlayerList
	pNum, ok := helper.StringInStrings(player, playerList)
	if !ok {
		return "", fmt.Errorf("could not find a player named %s", player)
	}

	a := command.ExtractNamedCommandArgs(args)
	concedePlayers := []string{}
	if len(a) == 0 && len(playerList) == 2 {
		// Two player, default to conceding to other player
		concedePlayers = []string{playerList[(pNum+1)%2]}
	} else {
		matchedPlayers := map[string]bool{}
		for _, p := range a {
			i, err := helper.MatchStringInStrings(p, playerList)
			if err != nil {
				return "", err
			}
			concedePlayer := playerList[i]
			if !matchedPlayers[concedePlayer] {
				matchedPlayers[concedePlayer] = true
				concedePlayers = append(concedePlayers, concedePlayer)
			}
		}
		if len(concedePlayers) == 0 {
			return "", errors.New("you must specify a player to concede to")
		}
	}
	c.gameModel.ConcedePlayers = concedePlayers
	c.gameModel.ConcedeVote = map[string]bool{
		player: true,
	}
	if len(concedePlayers) == 1 {
		// If it's a vote for a single winner, they automatically vote true
		c.gameModel.ConcedeVote[concedePlayers[0]] = true
	}

	if len(c.gameModel.RemainingConcedeVotePlayers()) == 0 {
		PassConcedeVote(c.gameModel, g)
	} else {
		g.GameLog().Add(log.NewPublicMessage(fmt.Sprintf(
			"%s called a vote to concede the game to %s",
			render.PlayerName(pNum, player),
			render.CommaList(render.PlayerNamesInPlayers(
				concedePlayers,
				playerList,
			)),
		)))
	}
	return "", nil
}

func (c ConcedeCommand) Usage(player string, context interface{}) string {
	return "{{b}}concede (## ##){{_b}} to concede to one or more other players"
}
