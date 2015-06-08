package email

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game"
	"github.com/Miniand/brdg.me/game/log"
	"github.com/Miniand/brdg.me/render"
	"github.com/Miniand/brdg.me/server/model"
)

func SendGame(
	g game.Playable,
	gm *model.GameModel,
	to []string,
	commands []command.Command,
	header string,
	initial bool,
) error {
	if header != "" {
		header += "\n\n"
	}
	if gm.IsFinished {
		winners := gm.Winners
		header += "The game is over"
		if len(winners) == 0 {
			header += ", it was a draw!"
		} else {
			header += ", the winners were: " + strings.Join(
				render.PlayerNamesInPlayers(winners, gm.PlayerList), ", ")
		}
	} else {
		header += "Current turn: " + strings.Join(render.PlayerNamesInPlayers(
			gm.WhoseTurn, gm.PlayerList), ", ")
	}
	commErrs := []string{}
	for _, p := range to {
		u, ok, err := model.FirstUserByEmail(p)
		if err != nil || ok && u.Unsubscribed {
			continue
		}
		pHeader := header
		rawOutput, err := g.RenderForPlayer(p)
		if err != nil {
			commErrs = append(commErrs, err.Error())
			continue
		}
		// Add log to header if needed
		messages := g.GameLog().NewMessagesFor(p)
		if len(messages) > 0 {
			pHeader += "\n\n{{b}}Since last time:{{_b}}\n" +
				log.RenderMessages(messages)
		}
		g.GameLog().MarkReadFor(p)
		// Add usages to header if needed
		usages := command.CommandUsages(p, g,
			command.AvailableCommands(p, g, commands))
		if len(usages) > 0 {
			pHeader += "\n\n{{b}}You can:{{_b}}\n" +
				render.CommandUsages(usages)
		}
		body := fmt.Sprintf(
			`%s

%s


You can also <a href="http://brdg.me/game.html?id=%s">continue playing this game live in your browser</a>.`,
			pHeader,
			rawOutput,
			gm.Id,
		)
		subject := fmt.Sprintf("%s (%s)", g.Name(), gm.Id)
		extraHeaders := []string{}
		messageId := gm.Id + "@brdg.me"
		if initial {
			// We create the base Message-ID
			extraHeaders = append(extraHeaders,
				fmt.Sprintf("Message-Id: <%s>", messageId))
		} else {
			// We create a unique Message-ID and set the In-Reply-To to original
			extraHeaders = append(extraHeaders,
				fmt.Sprintf("In-Reply-To: <%s>", messageId))
		}
		err = SendRichMail([]string{p}, subject, body, extraHeaders)
		if err != nil {
			commErrs = append(commErrs, err.Error())
			continue
		}
	}
	if len(commErrs) > 0 {
		return errors.New(strings.Join(commErrs, "\n"))
	}
	return nil
}
