package game

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/scommand"
)

var gameMut = map[string]*sync.Mutex{}

// Run commands on the game, email relevant people and handle action issues
func HandleCommandText(player, gameId, commandText string) error {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || ok && u.Unsubscribed || gameId == "" {
		commands := scommand.Commands(nil)
		output, err := command.CallInCommands(player, nil, commandText, commands)
		if err != nil {
			// Print help
			body := bytes.NewBufferString("")
			if err != command.NO_COMMAND_FOUND {
				body.WriteString(err.Error())
				body.WriteString("\n\n")
			}
			if output != "" {
				body.WriteString(output)
				body.WriteString("\n\n")
			}
			body.WriteString("Welcome to brdg.me!\n\n")
			body.WriteString(render.CommandUsages(
				command.CommandUsages(player, nil,
					command.AvailableCommands(player, nil, commands))))
			err = email.SendRichMail([]string{player}, "Welcome to brdg.me!",
				body.String(), []string{})
			if err != nil {
				return err
			}
		}
	} else {
		// Ensure games are loaded, modified and saved only one at a time.
		if gameMut[gameId] == nil {
			gameMut[gameId] = &sync.Mutex{}
		}
		gameMut[gameId].Lock()
		exitedBeforeSave := true
		defer func() {
			// This unlock logic will work until we defer for updating the
			// game, after which we want to unlock the mutex after update.
			if exitedBeforeSave {
				gameMut[gameId].Unlock()
			}
		}()
		var initialEliminated []string
		gm, err := model.LoadGame(gameId)
		if err != nil {
			return err
		}
		g, err := gm.ToGame()
		if err != nil {
			return err
		}
		exitedBeforeSave = false // Don't use previous deferred unlock.
		defer func() {
			if err := gm.UpdateState(g); err == nil {
				gm.Save()
			}
			gameMut[gameId].Unlock()
		}()
		alreadyFinished := g.IsFinished()
		commands := append(g.Commands(), scommand.Commands(gm)...)
		initialWhoseTurn := g.WhoseTurn()
		eliminator, isEliminator := g.(game.Eliminator)
		if isEliminator {
			initialEliminated = eliminator.EliminatedPlayerList()
		}
		output, err := command.CallInCommands(player, g, commandText, commands)
		header := ""
		if err != nil {
			header = err.Error()
		}
		if output != "" {
			if header != "" {
				header += "\n\n"
			}
			header += output
		}
		commErrs := []string{}
		commErr := communicate.Game(gm.Id, g, []string{player},
			commands, header, false)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if err != command.NO_COMMAND_FOUND {
			// Keep track who we've communicated to for if it's the end of the
			// game.
			communicatedTo := []string{player}
			// Email any players who now have a turn, or for ones who still have
			// a turn but there are new logs
			whoseTurnNow, remaining := WhoseTurnNow(g, initialWhoseTurn)
			commErr = communicate.Game(gm.Id, g, whoseTurnNow,
				append(g.Commands(), scommand.Commands(gm)...), "", false)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			communicatedTo = append(communicatedTo, whoseTurnNow...)
			whoseTurnNewLogs := []string{}
			for _, p := range remaining {
				if len(g.GameLog().NewMessagesFor(p)) > 0 {
					whoseTurnNewLogs = append(whoseTurnNewLogs, p)
				}
			}
			commErr = communicate.Game(gm.Id, g, whoseTurnNewLogs,
				append(g.Commands(), scommand.Commands(gm)...),
				"", false)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			communicatedTo = append(communicatedTo, whoseTurnNewLogs...)
			// Email any players who were eliminated this turn
			if isEliminator {
				newlyEliminated, _ := FindNewStringsInSlice(initialEliminated,
					eliminator.EliminatedPlayerList())
				commErr = communicate.Game(gm.Id, g, newlyEliminated,
					append(g.Commands(), scommand.Commands(gm)...),
					"You have been eliminated from the game.", false)
				if commErr != nil {
					commErrs = append(commErrs, commErr.Error())
				}
				communicatedTo = append(communicatedTo, newlyEliminated...)
			}
			uncommunicated, _ := FindNewStringsInSlice(communicatedTo,
				g.PlayerList())
			if len(uncommunicated) > 0 {
				if !alreadyFinished && g.IsFinished() {
					// If it's the end of the game and some people haven't been contacted
					commErr = communicate.Game(gm.Id, g, uncommunicated,
						append(g.Commands(), scommand.Commands(gm)...), "", false)
				} else {
					// We send updates to all remaining players via websockets so
					// they can update.
					communicate.GameUpdate(gm.Id, g, uncommunicated, "")
				}
			}
		}
		if len(commErrs) > 0 {
			return errors.New(strings.Join(commErrs, "\n"))
		}
	}
	return nil
}

func WhoseTurnNow(g game.Playable, initialWhoseTurn []string) ([]string,
	[]string) {
	return FindNewStringsInSlice(initialWhoseTurn, g.WhoseTurn())
}

func FindNewStringsInSlice(oldSlice, newSlice []string) (newStrings,
	remaining []string) {
	oldSliceMap := map[string]bool{}
	for _, s := range oldSlice {
		oldSliceMap[s] = true
	}
	for _, s := range newSlice {
		if !oldSliceMap[s] {
			newStrings = append(newStrings, s)
		} else {
			remaining = append(remaining, s)
		}
	}
	return
}
