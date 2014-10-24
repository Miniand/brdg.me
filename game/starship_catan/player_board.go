package starship_catan

import (
	"reflect"

	"github.com/Miniand/brdg.me/game/card"
)

type PlayerBoard struct {
	Player              int
	Resources           map[int]int
	Modules             map[int]int
	CompletedAdventures card.Deck
	Colonies            card.Deck
	TradingPosts        card.Deck
	DefeatedPirates     card.Deck
	FriendOfThePeople   bool
	HeroOfThePeople     bool
	LastSectors         []int
}

func NewPlayerBoard(player int) *PlayerBoard {
	pb := &PlayerBoard{
		Player: player,
		Resources: map[int]int{
			ResourceTrade:      2,
			ResourceScience:    1,
			ResourceAstro:      25,
			ResourceColonyShip: 1,
			ResourceTradeShip:  1,
			ResourceBooster:    2,
			ResourceCannon:     1,
		},
		Modules:             map[int]int{},
		CompletedAdventures: card.Deck{},
		Colonies:            card.Deck{},
		TradingPosts:        card.Deck{},
		DefeatedPirates:     card.Deck{},
		LastSectors:         []int{},
	}
	pb.Colonies = pb.Colonies.Push(StartingCards()[player])
	return pb
}

func (b *PlayerBoard) Actions() int {
	return 2 + b.Modules[ModuleCommand]
}

func (b *PlayerBoard) Ships() int {
	return b.Resources[ResourceTradeShip] + b.Resources[ResourceColonyShip]
}

func (b *PlayerBoard) CanBuildShip() bool {
	return b.Ships() < 2
}

func (b *PlayerBoard) CanBuildBooster() bool {
	return b.Resources[ResourceBooster] < 6
}

func (b *PlayerBoard) CanBuildCannon() bool {
	return b.Resources[ResourceCannon] < 6
}

func (b *PlayerBoard) CanBuild() bool {
	return b.CanBuildShip() || b.CanBuildBooster() || b.CanBuildCannon()
}

func (b *PlayerBoard) BoosterTransaction() Transaction {
	t := Transaction{
		ResourceFuel:    -2,
		ResourceBooster: 1,
	}
	if b.Resources[ResourceBooster] >= 3 {
		t[ResourceScience] = -1
	}
	return t
}

func (b *PlayerBoard) CannonTransaction() Transaction {
	t := Transaction{
		ResourceCarbon: -2,
		ResourceCannon: 1,
	}
	if b.Resources[ResourceBooster] >= 3 {
		t[ResourceScience] = -1
	}
	return t
}

func (b *PlayerBoard) CanAfford(t Transaction) bool {
	return b.CanFit(t.Lose())
}

func (b *PlayerBoard) GoodsLimit() int {
	return 2 + b.Modules[ModuleLogistics]
}

func (b *PlayerBoard) CanFit(t Transaction) bool {
	return reflect.DeepEqual(t, b.FitTransaction(t))
}

func (b *PlayerBoard) FitTransaction(t Transaction) Transaction {
	fit := Transaction{}
	for r, v := range t {
		fit[r] = v
	}
	shipTotal := b.Ships() + fit[ResourceColonyShip] + fit[ResourceTradeShip]
	for shipTotal > 2 {
		if fit[ResourceColonyShip] > 0 {
			fit[ResourceColonyShip] -= 1
		} else {
			fit[ResourceTradeShip] -= 1
		}
		shipTotal -= 1
	}
	goodsLimit := b.GoodsLimit()
	for r, v := range fit {
		// Ensure none too large
		switch {
		case r == ResourceScience:
			if b.Resources[r]+v > 4 {
				fit[r] = 4 - b.Resources[r]
			}
		case r >= ResourceFood && r <= ResourceOre || r == ResourceTrade:
			if b.Resources[r]+v > goodsLimit {
				fit[r] = goodsLimit - b.Resources[r]
			}
		case r == ResourceBooster || r == ResourceCannon:
			if b.Resources[r]+v > 6 {
				fit[r] = 6 - b.Resources[r]
			}
		}
		// Ensure we don't go below 0
		if b.Resources[r]+fit[r] < 0 {
			fit[r] = -b.Resources[r]
		}
	}
	return fit.TrimEmpty()
}

func (b *PlayerBoard) Transact(t Transaction) Transaction {
	fitted := b.FitTransaction(t)
	for r, v := range fitted {
		b.Resources[r] += v
	}
	return fitted
}

type PlayerTradingPrices struct {
	Buy, Sell int
}

func (b *PlayerBoard) TradingPostPrices() map[int]*PlayerTradingPrices {
	prices := map[int]*PlayerTradingPrices{}
	for _, c := range b.TradingPosts {
		if tp, ok := c.(TradeCard); ok {
			for _, r := range tp.Resources {
				if prices[r] == nil {
					prices[r] = &PlayerTradingPrices{}
				}
				if tp.Direction != TradeDirSell &&
					(prices[r].Buy == 0 || tp.Price < prices[r].Buy) {
					prices[r].Buy = tp.Price
				}
				if tp.Direction != TradeDirBuy &&
					(prices[r].Sell == 0 || tp.Price > prices[r].Sell) {
					prices[r].Sell = tp.Price
				}
			}
		}
	}
	return prices
}
