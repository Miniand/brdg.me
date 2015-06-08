package seven_wonders

import (
	"testing"

	"github.com/Miniand/brdg.me/game/card"
	"github.com/stretchr/testify/assert"
)

func TestCardMimicGuild(t *testing.T) {
	g := &Game{}
	assert.NoError(t, g.Start(players))

	g.Cards[Mick] = card.Deck{
		Cards[WonderStageOlympiaB3],
	}
	g.Cards[Steve] = card.Deck{
		Cards[CardBuildersGuild], // Should be worth 1 thanks to Mick's wonder
	}
	g.Cards[Greg] = card.Deck{
		Cards[CardWorkersGuild], // Should be worth 0
	}

	assert.Equal(t, 2, g.PlayerResourceCount(Mick, VP))
}
