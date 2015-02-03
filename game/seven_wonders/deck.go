package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

func DeckAge1() card.Deck {
	d := card.Deck{}
	for _, c := range []card.Deck{
		// Raw goods
		NewCardGoodRaw(CardLumberYard, nil, []int{GoodWood}, 1, 3, 4),
		NewCardGoodRaw(CardStonePit, nil, []int{GoodStone}, 1, 3, 5),
		NewCardGoodRaw(CardClayPool, nil, []int{GoodClay}, 1, 3, 5),
		NewCardGoodRaw(CardOreVein, nil, []int{GoodOre}, 1, 3, 4),
		NewCardGoodRaw(CardTreeFarm, Cost{
			GoodCoin: 1,
		}, []int{GoodWood, GoodClay}, 1, 6),
		NewCardGoodRaw(CardExcavation, Cost{
			GoodCoin: 1,
		}, []int{GoodStone, GoodClay}, 1, 4),
		NewCardGoodRaw(CardClayPit, Cost{
			GoodCoin: 1,
		}, []int{GoodClay, GoodOre}, 1, 3),
		NewCardGoodRaw(CardTimberYard, Cost{
			GoodCoin: 1,
		}, []int{GoodStone, GoodWood}, 1, 3),
		NewCardGoodRaw(CardForestCave, Cost{
			GoodCoin: 1,
		}, []int{GoodWood, GoodOre}, 1, 5),
		NewCardGoodRaw(CardMine, Cost{
			GoodCoin: 1,
		}, []int{GoodStone, GoodOre}, 1, 6),

		// Manufactured goods
		NewCardGoodManufactured(CardLoom, nil, []int{GoodTextile}, 1, 3, 6),
		NewCardGoodManufactured(CardGlassworks, nil, []int{GoodGlass}, 1, 3, 6),
		NewCardGoodManufactured(CardPress, nil, []int{GoodPapyrus}, 1, 3, 6),

		// Civilian structures
		NewCardCivilian(CardPawnshop, nil, 3, nil, nil, 4, 7),
		NewCardCivilian(CardBaths, Cost{
			GoodStone: 1,
		}, 3, nil, []string{CardAqueduct}, 3, 7),
		NewCardCivilian(CardAltar, nil, 2, nil, []string{CardTemple}, 3, 5),
		NewCardCivilian(CardTheater, nil, 2, nil, []string{CardStatue}, 3, 6),

		// Commercial structures
		NewCardCommercialTavern(4, 5, 7),
		NewCardCommercialTrade(
			CardEastTradingPost,
			nil,
			[]int{DirRight},
			RawGoods,
			nil,
			[]string{CardForum},
			3, 7,
		),
		NewCardCommercialTrade(
			CardWestTradingPost,
			nil,
			[]int{DirLeft},
			RawGoods,
			nil,
			[]string{CardForum},
			3, 7,
		),
		NewCardCommercialTrade(
			CardMarketplace,
			nil,
			[]int{DirLeft, DirRight},
			ManufacturedGoods,
			nil,
			[]string{CardCaravansery},
			3, 6,
		),

		// Military structures
		NewCardMilitary(CardStockade, Cost{
			GoodWood: 1,
		}, 1, nil, nil, 3, 7),
		NewCardMilitary(CardBarracks, Cost{
			GoodOre: 1,
		}, 1, nil, nil, 3, 5),
		NewCardMilitary(CardGuardTower, Cost{
			GoodClay: 1,
		}, 1, nil, nil, 3, 4),

		// Scientific structures
		NewCardScience(CardApothecary, Cost{
			GoodTextile: 1,
		}, FieldMathematics, nil, []string{
			CardStables,
			CardDispensary,
		}, 3, 5),
		NewCardScience(CardWorkshop, Cost{
			GoodGlass: 1,
		}, FieldEngineering, nil, []string{
			CardLaboratory,
			CardArcheryRange,
		}, 3, 7),
		NewCardScience(CardScriptorium, Cost{
			GoodPapyrus: 1,
		}, FieldTheology, nil, []string{
			CardCourthouse,
			CardLibrary,
		}, 3, 4),
	} {
		d = d.PushMany(c)
	}
	return d
}
