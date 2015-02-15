package seven_wonders

import (
	"github.com/Miniand/brdg.me/game/card"
	"github.com/Miniand/brdg.me/game/cost"
)

var Cards = map[string]Carder{
	// Age 1

	// Raw goods
	CardLumberYard: NewCardGoodRaw(CardLumberYard, nil, []cost.Cost{
		{GoodWood: 1},
	}),
	CardStonePit: NewCardGoodRaw(CardStonePit, nil, []cost.Cost{
		{GoodStone: 1},
	}),
	CardClayPool: NewCardGoodRaw(CardClayPool, nil, []cost.Cost{
		{GoodClay: 1},
	}),
	CardOreVein: NewCardGoodRaw(CardOreVein, nil, []cost.Cost{
		{GoodOre: 1},
	}),
	CardTreeFarm: NewCardGoodRaw(CardTreeFarm, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodWood: 1},
		{GoodClay: 1},
	}),
	CardExcavation: NewCardGoodRaw(CardExcavation, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodStone: 1},
		{GoodClay: 1},
	}),
	CardClayPit: NewCardGoodRaw(CardClayPit, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodClay: 1},
		{GoodOre: 1},
	}),
	CardTimberYard: NewCardGoodRaw(CardTimberYard, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodStone: 1},
		{GoodWood: 1},
	}),
	CardForestCave: NewCardGoodRaw(CardForestCave, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodWood: 1},
		{GoodOre: 1},
	}),
	CardMine: NewCardGoodRaw(CardMine, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodStone: 1},
		{GoodOre: 1},
	}),

	// Manufactured goods
	CardLoom: NewCardGoodManufactured(CardLoom, nil, []cost.Cost{
		{GoodTextile: 1},
	}),
	CardGlassworks: NewCardGoodManufactured(CardGlassworks, nil, []cost.Cost{
		{GoodGlass: 1},
	}),
	CardPress: NewCardGoodManufactured(CardPress, nil, []cost.Cost{
		{GoodPapyrus: 1},
	}),

	// Civilian structures
	CardPawnshop: NewCardCivilian(CardPawnshop, nil, 3, nil, nil),
	CardBaths: NewCardCivilian(CardBaths, cost.Cost{
		GoodStone: 1,
	}, 3, nil, []string{CardAqueduct}),
	CardAltar:   NewCardCivilian(CardAltar, nil, 2, nil, []string{CardTemple}),
	CardTheater: NewCardCivilian(CardTheater, nil, 2, nil, []string{CardStatue}),

	// Commercial structures
	CardTavern: NewCardCommercialTavern(),
	CardEastTradingPost: NewCardCommercialTrade(
		CardEastTradingPost,
		nil,
		[]int{DirRight},
		RawGoods,
		nil,
		[]string{CardForum},
	),
	CardWestTradingPost: NewCardCommercialTrade(
		CardWestTradingPost,
		nil,
		[]int{DirLeft},
		RawGoods,
		nil,
		[]string{CardForum},
	),
	CardMarketplace: NewCardCommercialTrade(
		CardMarketplace,
		nil,
		[]int{DirLeft, DirRight},
		ManufacturedGoods,
		nil,
		[]string{CardCaravansery},
	),

	// Military structures
	CardStockade: NewCardMilitary(CardStockade, cost.Cost{
		GoodWood: 1,
	}, 1, nil, nil),
	CardBarracks: NewCardMilitary(CardBarracks, cost.Cost{
		GoodOre: 1,
	}, 1, nil, nil),
	CardGuardTower: NewCardMilitary(CardGuardTower, cost.Cost{
		GoodClay: 1,
	}, 1, nil, nil),

	// Scientific structures
	CardApothecary: NewCardScience(CardApothecary, cost.Cost{
		GoodTextile: 1,
	}, FieldMathematics, nil, []string{
		CardStables,
		CardDispensary,
	}),
	CardWorkshop: NewCardScience(CardWorkshop, cost.Cost{
		GoodGlass: 1,
	}, FieldEngineering, nil, []string{
		CardLaboratory,
		CardArcheryRange,
	}),
	CardScriptorium: NewCardScience(CardScriptorium, cost.Cost{
		GoodPapyrus: 1,
	}, FieldTheology, nil, []string{
		CardCourthouse,
		CardLibrary,
	}),

	// Age 2

	// Raw goods
	CardSawmill: NewCardGoodRaw(CardSawmill, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodWood: 2},
	}),
	CardQuarry: NewCardGoodRaw(CardQuarry, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodStone: 2},
	}),
	CardBrickyard: NewCardGoodRaw(CardBrickyard, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodClay: 2},
	}),
	CardFoundry: NewCardGoodRaw(CardFoundry, cost.Cost{
		GoodCoin: 1,
	}, []cost.Cost{
		{GoodOre: 2},
	}),

	// Civilian structures
	CardAqueduct: NewCardCivilian(CardAqueduct, cost.Cost{
		GoodStone: 3,
	}, 5, []string{CardBaths}, nil),
	CardTemple: NewCardCivilian(CardTemple, cost.Cost{
		GoodWood:  1,
		GoodClay:  1,
		GoodGlass: 1,
	}, 3, []string{CardAltar}, []string{CardPantheon}),
	CardStatue: NewCardCivilian(CardStatue, cost.Cost{
		GoodOre:  2,
		GoodWood: 1,
	}, 4, []string{CardTheater}, []string{CardGardens}),
	CardCourthouse: NewCardCivilian(CardCourthouse, cost.Cost{
		GoodClay:    2,
		GoodTextile: 1,
	}, 4, []string{CardScriptorium}, nil),

	// Commercial structures
	CardForum: NewCardGoodCommercial(CardForum, cost.Cost{
		GoodClay: 2,
	}, SliceToCost(ManufacturedGoods), []string{
		CardEastTradingPost,
		CardWestTradingPost,
	}, []string{CardHaven}),
	CardCaravansery: NewCardGoodCommercial(CardCaravansery, cost.Cost{
		GoodWood: 2,
	}, SliceToCost(RawGoods), []string{CardMarketplace}, []string{CardLighthouse}),
	CardVineyard: NewCardBonus(CardVineyard, CardKindCommercial, nil, []int{
		CardKindRaw,
	}, DirAll, 0, 1, nil, nil),
	CardBazar: NewCardBonus(CardBazar, CardKindCommercial, nil, []int{
		CardKindManufactured,
	}, DirAll, 0, 2, nil, nil),

	// Military structures
	CardWalls: NewCardMilitary(CardWalls, cost.Cost{
		GoodStone: 3,
	}, 2, nil, []string{CardFortifications}),
	CardTrainingGround: NewCardMilitary(CardTrainingGround, cost.Cost{
		GoodOre:  2,
		GoodWood: 1,
	}, 2, nil, []string{CardCircus}),
	CardStables: NewCardMilitary(CardStables, cost.Cost{
		GoodClay: 1,
		GoodWood: 1,
		GoodOre:  1,
	}, 2, []string{CardApothecary}, nil),
	CardArcheryRange: NewCardMilitary(CardArcheryRange, cost.Cost{
		GoodWood: 2,
		GoodOre:  1,
	}, 2, []string{CardWorkshop}, nil),

	// Scientific structures
	CardDispensary: NewCardScience(CardDispensary, cost.Cost{
		GoodOre:   2,
		GoodGlass: 1,
	}, FieldMathematics, []string{CardApothecary}, []string{
		CardLodge,
		CardArena,
	}),
	CardLaboratory: NewCardScience(CardLaboratory, cost.Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
	}, FieldEngineering, []string{CardWorkshop}, []string{
		CardObservatory,
		CardSiegeWorkshop,
	}),
	CardLibrary: NewCardScience(CardLibrary, cost.Cost{
		GoodStone:   2,
		GoodTextile: 1,
	}, FieldTheology, []string{CardScriptorium}, []string{
		CardSenate,
		CardUniversity,
	}),
	CardSchool: NewCardScience(CardSchool, cost.Cost{
		GoodWood:    1,
		GoodPapyrus: 1,
	}, FieldTheology, nil, []string{
		CardAcademy,
		CardStudy,
	}),

	// Age 3

	// Civilian structures
	CardPantheon: NewCardCivilian(CardPantheon, cost.Cost{
		GoodClay:    2,
		GoodOre:     1,
		GoodGlass:   1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, 7, []string{CardTemple}, nil),
	CardGardens: NewCardCivilian(CardGardens, cost.Cost{
		GoodClay: 2,
		GoodWood: 1,
	}, 5, []string{CardStatue}, nil),
	CardTownHall: NewCardCivilian(CardTownHall, cost.Cost{
		GoodStone: 2,
		GoodOre:   1,
		GoodGlass: 1,
	}, 6, nil, nil),
	CardPalace: NewCardCivilian(CardPalace, cost.Cost{
		GoodStone:   1,
		GoodOre:     1,
		GoodWood:    1,
		GoodClay:    1,
		GoodGlass:   1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, 8, nil, nil),
	CardSenate: NewCardCivilian(CardSenate, cost.Cost{
		GoodWood:  2,
		GoodStone: 1,
		GoodOre:   1,
	}, 6, []string{CardLibrary}, nil),

	// Commercial structures
	CardHaven: NewCardBonus(CardHaven, CardKindCommercial, cost.Cost{
		GoodWood:    1,
		GoodOre:     1,
		GoodTextile: 1,
	}, []int{CardKindRaw}, DirSelf, 1, 1, []string{CardForum}, nil),
	CardLighthouse: NewCardBonus(CardLighthouse, CardKindCommercial, cost.Cost{
		GoodStone: 1,
		GoodGlass: 1,
	}, []int{CardKindCommercial}, DirSelf, 1, 1, []string{
		CardCaravansery,
	}, nil),
	CardChamberOfCommerce: NewCardBonus(CardChamberOfCommerce, CardKindCommercial, cost.Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
	}, []int{CardKindManufactured}, DirSelf, 2, 2, nil, nil),
	CardArena: NewCardBonus(CardArena, CardKindCommercial, cost.Cost{
		GoodStone: 2,
		GoodOre:   1,
	}, []int{WonderStage}, DirSelf, 1, 3, []string{
		CardDispensary,
	}, nil),

	// Military structures
	CardFortifications: NewCardMilitary(CardFortifications, cost.Cost{
		GoodOre:   3,
		GoodStone: 1,
	}, 3, []string{CardWalls}, nil),
	CardCircus: NewCardMilitary(CardCircus, cost.Cost{
		GoodStone: 3,
		GoodOre:   1,
	}, 3, []string{CardTrainingGround}, nil),
	CardArsenal: NewCardMilitary(CardArsenal, cost.Cost{
		GoodWood:    2,
		GoodOre:     1,
		GoodTextile: 1,
	}, 3, nil, nil),
	CardSiegeWorkshop: NewCardMilitary(CardSiegeWorkshop, cost.Cost{
		GoodClay: 3,
		GoodWood: 1,
	}, 3, []string{CardLaboratory}, nil),

	// Scientific structures
	CardLodge: NewCardScience(CardLodge, cost.Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, FieldMathematics, []string{CardDispensary}, nil),
	CardObservatory: NewCardScience(CardObservatory, cost.Cost{
		GoodOre:     2,
		GoodGlass:   1,
		GoodTextile: 1,
	}, FieldEngineering, []string{CardLaboratory}, nil),
	CardUniversity: NewCardScience(CardUniversity, cost.Cost{
		GoodWood:    2,
		GoodPapyrus: 1,
		GoodGlass:   1,
	}, FieldTheology, []string{CardLibrary}, nil),
	CardAcademy: NewCardScience(CardAcademy, cost.Cost{
		GoodStone: 3,
		GoodGlass: 1,
	}, FieldMathematics, []string{CardSchool}, nil),
	CardStudy: NewCardScience(CardStudy, cost.Cost{
		GoodWood:    1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, FieldEngineering, []string{CardSchool}, nil),

	// Guilds

	CardWorkersGuild: NewCardBonus(CardWorkersGuild, CardKindGuild, cost.Cost{
		GoodOre:   2,
		GoodClay:  1,
		GoodStone: 1,
		GoodWood:  1,
	}, []int{CardKindRaw}, DirNeighbours, 1, 0, nil, nil),
	CardCraftsmensGuild: NewCardBonus(CardCraftsmensGuild, CardKindGuild, cost.Cost{
		GoodOre:   2,
		GoodStone: 2,
	}, []int{CardKindManufactured}, DirNeighbours, 2, 0, nil, nil),
	CardTradersGuild: NewCardBonus(CardTradersGuild, CardKindGuild, cost.Cost{
		GoodGlass:   1,
		GoodTextile: 1,
		GoodPapyrus: 1,
	}, []int{CardKindCommercial}, DirNeighbours, 1, 0, nil, nil),
	CardPhilosophersGuild: NewCardBonus(CardPhilosophersGuild, CardKindGuild, cost.Cost{
		GoodClay:    3,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, []int{CardKindScientific}, DirNeighbours, 1, 0, nil, nil),
	CardSpiesGuild: NewCardBonus(CardSpiesGuild, CardKindGuild, cost.Cost{
		GoodClay:  3,
		GoodGlass: 1,
	}, []int{CardKindMilitary}, DirNeighbours, 1, 0, nil, nil),
	CardStrategistsGuild: NewCardBonus(CardStrategistsGuild, CardKindGuild, cost.Cost{
		GoodOre:     2,
		GoodStone:   1,
		GoodTextile: 1,
	}, []int{TokenDefeat}, DirNeighbours, 1, 0, nil, nil),
	CardShipownersGuild: NewCardBonus(CardShipownersGuild, CardKindGuild, cost.Cost{
		GoodWood:    3,
		GoodGlass:   1,
		GoodPapyrus: 1,
	}, []int{
		CardKindRaw,
		CardKindManufactured,
		CardKindGuild,
	}, DirSelf, 1, 0, nil, nil),
	CardScientistsGuild: CardScience{
		Card{
			CardScientistsGuild,
			CardKindGuild,
			cost.Cost{
				GoodWood:    2,
				GoodOre:     2,
				GoodPapyrus: 1,
			},
			nil,
			nil,
		},
		AllFields,
	},
	CardMagistratesGuild: NewCardBonus(CardMagistratesGuild, CardKindGuild, cost.Cost{
		GoodWood:    3,
		GoodStone:   1,
		GoodTextile: 1,
	}, []int{CardKindCivilian}, DirNeighbours, 1, 0, nil, nil),
	CardBuildersGuild: NewCardBonus(CardBuildersGuild, CardKindGuild, cost.Cost{
		GoodStone: 2,
		GoodClay:  2,
		GoodGlass: 1,
	}, []int{WonderStage}, DirAll, 1, 0, nil, nil),

	// Wonder stages

	WonderStageRhodesA1: NewCardWonderVP(
		WonderStageRhodesA1,
		cost.Cost{GoodWood: 2},
		3,
	),
	WonderStageRhodesA2: CardMilitary{
		NewCard(
			WonderStageRhodesA2,
			CardKindWonder,
			cost.Cost{GoodClay: 3},
			nil,
			nil,
		),
		2,
	},
	WonderStageRhodesA3: NewCardWonderVP(
		WonderStageRhodesA3,
		cost.Cost{GoodOre: 4},
		7,
	),

	WonderStageRhodesB1: CardMulti{
		NewCard(
			WonderStageRhodesB1,
			CardKindWonder,
			cost.Cost{GoodStone: 3},
			nil,
			nil,
		),
		cost.Cost{
			AttackStrength: 1,
			VP:             3,
			GoodCoin:       3,
		},
	},
	WonderStageRhodesB2: CardMulti{
		NewCard(
			WonderStageRhodesB2,
			CardKindWonder,
			cost.Cost{GoodOre: 4},
			nil,
			nil,
		),
		cost.Cost{
			AttackStrength: 1,
			VP:             4,
			GoodCoin:       4,
		},
	},
}

func DeckForPlayers(cards []CardForPlayers, players int) card.Deck {
	d := card.Deck{}
	for _, c := range cards {
		for _, p := range c.Players {
			if p <= players {
				d = d.Push(Cards[c.Card])
			}
		}
	}
	return d
}

func DeckAge1(players int) card.Deck {
	return DeckForPlayers([]CardForPlayers{
		// Raw goods
		{CardLumberYard, []int{3, 4}},
		{CardStonePit, []int{3, 5}},
		{CardClayPool, []int{3, 5}},
		{CardOreVein, []int{3, 4}},
		{CardTreeFarm, []int{6}},
		{CardExcavation, []int{4}},
		{CardClayPit, []int{3}},
		{CardTimberYard, []int{3}},
		{CardForestCave, []int{5}},
		{CardMine, []int{6}},

		// Manufactured goods
		{CardLoom, []int{3, 6}},
		{CardGlassworks, []int{3, 6}},
		{CardPress, []int{3, 6}},

		// Civilian structures
		{CardPawnshop, []int{4, 7}},
		{CardBaths, []int{3, 7}},
		{CardAltar, []int{3, 5}},
		{CardTheater, []int{3, 6}},

		// Commercial structures
		{CardTavern, []int{4, 5, 7}},
		{CardEastTradingPost, []int{3, 7}},
		{CardWestTradingPost, []int{3, 7}},
		{CardMarketplace, []int{3, 6}},

		// Military structures
		{CardStockade, []int{3, 7}},
		{CardBarracks, []int{3, 5}},
		{CardGuardTower, []int{3, 4}},

		// Scientific structures
		{CardApothecary, []int{3, 5}},
		{CardWorkshop, []int{3, 7}},
		{CardScriptorium, []int{3, 4}},
	}, players)
}

func DeckAge2(players int) card.Deck {
	return DeckForPlayers([]CardForPlayers{
		// Raw goods
		{CardSawmill, []int{3, 4}},
		{CardQuarry, []int{3, 4}},
		{CardBrickyard, []int{3, 4}},
		{CardFoundry, []int{3, 4}},

		// Manufactured goods
		{CardLoom, []int{3, 5}},
		{CardGlassworks, []int{3, 5}},
		{CardPress, []int{3, 5}},

		// Civilian structures
		{CardAqueduct, []int{3, 7}},
		{CardTemple, []int{3, 6}},
		{CardStatue, []int{3, 7}},
		{CardCourthouse, []int{3, 5}},

		// Commercial structures
		{CardForum, []int{3, 6, 7}},
		{CardCaravansery, []int{3, 5, 6}},
		{CardVineyard, []int{3, 6}},
		{CardBazar, []int{4, 7}},

		// Military structures
		{CardWalls, []int{3, 7}},
		{CardTrainingGround, []int{4, 6, 7}},
		{CardStables, []int{3, 5}},
		{CardArcheryRange, []int{3, 6}},

		// Scientific structures
		{CardDispensary, []int{3, 4}},
		{CardLaboratory, []int{3, 5}},
		{CardLibrary, []int{3, 6}},
		{CardSchool, []int{3, 7}},
	}, players)
}

func DeckAge3(players int) card.Deck {
	return DeckForPlayers([]CardForPlayers{
		// Civilian structures
		{CardPantheon, []int{3, 6}},
		{CardGardens, []int{3, 4}},
		{CardTownHall, []int{3, 5, 6}},
		{CardPalace, []int{3, 7}},
		{CardSenate, []int{3, 5}},

		// Commercial structures
		{CardHaven, []int{3, 4}},
		{CardLighthouse, []int{3, 6}},
		{CardChamberOfCommerce, []int{4, 6}},
		{CardArena, []int{3, 5, 7}},

		// Military structures
		{CardFortifications, []int{3, 7}},
		{CardCircus, []int{4, 5, 6}},
		{CardArsenal, []int{3, 4, 7}},
		{CardSiegeWorkshop, []int{3, 5}},

		// Scientific structures
		{CardLodge, []int{3, 6}},
		{CardObservatory, []int{3, 7}},
		{CardUniversity, []int{3, 4}},
		{CardAcademy, []int{3, 7}},
		{CardStudy, []int{3, 5}},
	}, players).PushMany(DeckGuild.Shuffle()[:players+2])
}

var DeckGuild = card.Deck{
	Cards[CardWorkersGuild],
	Cards[CardCraftsmensGuild],
	Cards[CardTradersGuild],
	Cards[CardPhilosophersGuild],
	Cards[CardSpiesGuild],
	Cards[CardStrategistsGuild],
	Cards[CardShipownersGuild],
	Cards[CardScientistsGuild],
	Cards[CardMagistratesGuild],
	Cards[CardBuildersGuild],
}
