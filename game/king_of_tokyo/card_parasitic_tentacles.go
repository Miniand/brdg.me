package king_of_tokyo

import "fmt"

type CardParasiticTentacles struct{}

func (c CardParasiticTentacles) Name() string {
	return "Parasitic Tentacles"
}

func (c CardParasiticTentacles) Description() string {
	return fmt.Sprintf(
		"{{b}}You can purchase cards from other monsters.{{_b}} Pay them the %s cost.",
		EnergySymbol,
	)
}

func (c CardParasiticTentacles) Cost() int {
	return 4
}

func (c CardParasiticTentacles) Kind() int {
	return CardKindKeep
}

func (c CardParasiticTentacles) ModifyBuyable(
	game *Game,
	player int,
	buyable []BuyableCard,
) []BuyableCard {
	for p := range game.Players {
		if p == player {
			continue
		}
		for _, c := range game.Boards[p].Cards {
			buyable = append(buyable, BuyableCard{
				Card:       c,
				From:       BuyFromPlayer,
				FromPlayer: p,
			})
		}
	}
	return buyable
}
