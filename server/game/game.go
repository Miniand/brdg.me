package game

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/communicate"
	"github.com/Miniand/brdg.me/server/email"
	"github.com/Miniand/brdg.me/server/model"
	"github.com/Miniand/brdg.me/server/scommand"
)

const (
	MsgTypeElimitate = "eliminate"
	MsgTypeFinish    = "finish"
	MsgTypeSuccess   = "success"
	MsgTypeError     = "error"
	MsgTypeYourTurn  = "yourTurn"
	MsgTypeNewLogs   = "newLogs"
	MsgTypeUpdate    = "update"
)

var gameMut = map[string]*sync.Mutex{}

// Run commands on the game, email relevant people and handle action issues
func HandleCommandText(player, gameId, commandText string) error {
	u, ok, err := model.FirstUserByEmail(player)
	if err != nil || ok && u.Unsubscribed || gameId == "" {
		commands := scommand.Commands(player, nil)
		output, err := command.CallInCommands(player, nil, commandText, commands)
		if err != nil || output != "" {
			// Print help
			body := bytes.NewBufferString("")
			if err != nil && err != command.ErrNoCommandFound {
				body.WriteString(err.Error())
				body.WriteString("\n\n")
			}
			if output != "" {
				body.WriteString(output)
				body.WriteString("\n\n")
			}
			body.WriteString("Welcome to brdg.me!\n\n")
			body.WriteString(render.CommandUsages(
				command.CommandUsages(player, nil, commands)))
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
		alreadyFinished := gm.IsFinished
		commands := scommand.CommandsForGame(player, gm, g)
		initialWhoseTurn := gm.WhoseTurn
		initialEliminated = gm.EliminatedPlayerList
		msgType := MsgTypeSuccess
		output, err := command.CallInCommandsPostHook(
			player,
			g,
			commandText,
			commands,
			func() error {
				return gm.UpdateState(g)
			},
		)
		header := ""
		if err != nil {
			msgType = MsgTypeError
			header = err.Error()
		}
		// Save now so websocket updates don't update old data
		if err := gm.UpdateState(g); err != nil {
			return err
		}
		if err := gm.Save(); err != nil {
			return err
		}
		if output != "" {
			if header != "" {
				header += "\n\n"
			}
			header += output
		}
		commErrs := []string{}
		commErr := communicate.Game(
			g,
			gm,
			[]string{player},
			commands,
			header,
			msgType,
			false,
		)
		if commErr != nil {
			commErrs = append(commErrs, commErr.Error())
		}
		if err != command.ErrNoCommandFound {
			// Keep track who we've communicated to for if it's the end of the
			// game.
			communicatedTo := []string{player}
			// Email any players who now have a turn, or for ones who still have
			// a turn but there are new logs
			whoseTurnNow, remaining := FindNewStringsInSlice(
				initialWhoseTurn, gm.WhoseTurn)
			commErr = communicate.Game(
				g,
				gm,
				whoseTurnNow,
				scommand.CommandsForGame(player, gm, g),
				"",
				MsgTypeYourTurn,
				false,
			)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			communicatedTo = append(communicatedTo, whoseTurnNow...)
			whoseTurnNewLogs := []string{}
			uncommunicated, _ := FindNewStringsInSlice(communicatedTo, remaining)
			for _, p := range uncommunicated {
				if len(g.GameLog().NewMessagesFor(p)) > 0 {
					whoseTurnNewLogs = append(whoseTurnNewLogs, p)
				}
			}
			commErr = communicate.Game(
				g,
				gm,
				whoseTurnNewLogs,
				scommand.CommandsForGame(player, gm, g),
				"",
				MsgTypeNewLogs,
				false,
			)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			communicatedTo = append(communicatedTo, whoseTurnNewLogs...)
			// Email any players who were eliminated this turn
			newlyEliminated, _ := FindNewStringsInSlice(initialEliminated,
				gm.EliminatedPlayerList)
			commErr = communicate.Game(
				g,
				gm,
				newlyEliminated,
				scommand.CommandsForGame(player, gm, g),
				"You have been eliminated from the game.",
				MsgTypeElimitate,
				false,
			)
			if commErr != nil {
				commErrs = append(commErrs, commErr.Error())
			}
			communicatedTo = append(communicatedTo, newlyEliminated...)
			uncommunicated, _ = FindNewStringsInSlice(communicatedTo,
				gm.PlayerList)
			if len(uncommunicated) > 0 {
				if !alreadyFinished && gm.IsFinished {
					// If it's the end of the game and some people haven't been contacted
					commErr = communicate.Game(
						g,
						gm,
						uncommunicated,
						scommand.CommandsForGame(player, gm, g),
						"",
						MsgTypeFinish,
						false,
					)
				} else {
					// We send updates to all remaining players via websockets so
					// they can update.
					communicate.GameUpdate(g, gm, uncommunicated, "", MsgTypeUpdate)
				}
			}
		}
		if len(commErrs) > 0 {
			return errors.New(strings.Join(commErrs, "\n"))
		}
	}
	return nil
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
