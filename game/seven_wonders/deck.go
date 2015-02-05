package seven_wonders

import "github.com/Miniand/brdg.me/game/card"

var Cards = map[string]Carder{
	// Age 1

	// Raw goods
	CardLumberYard: NewCardGoodRaw(CardLumberYard, nil, []int{GoodWood}, 1),
	CardStonePit:   NewCardGoodRaw(CardStonePit, nil, []int{GoodStone}, 1),
	CardClayPool:   NewCardGoodRaw(CardClayPool, nil, []int{GoodClay}, 1),
	CardOreVein:    NewCardGoodRaw(CardOreVein, nil, []int{GoodOre}, 1),
	CardTreeFarm: NewCardGoodRaw(CardTreeFarm, Cost{
		GoodCoin: 1,
	}, []int{GoodWood, GoodClay}, 1),
	CardExcavation: NewCardGoodRaw(CardExcavation, Cost{
		GoodCoin: 1,
	}, []int{GoodStone, GoodClay}, 1),
	CardClayPit: NewCardGoodRaw(CardClayPit, Cost{
		GoodCoin: 1,
	}, []int{GoodClay, GoodOre}, 1),
	CardTimberYard: NewCardGoodRaw(CardTimberYard, Cost{
		GoodCoin: 1,
	}, []int{GoodStone, GoodWood}, 1),
	CardForestCave: NewCardGoodRaw(CardForestCave, Cost{
		GoodCoin: 1,
	}, []int{GoodWood, GoodOre}, 1),
	CardMine: NewCardGoodRaw(CardMine, Cost{
		GoodCoin: 1,
	}, []int{GoodStone, GoodOre}, 1),

	// Manufactured goods
	CardLoom:       NewCardGoodManufactured(CardLoom, nil, []int{GoodTextile}, 1),
	CardGlassworks: NewCardGoodManufactured(CardGlassworks, nil, []int{GoodGlass}, 1),
	CardPress:      NewCardGoodManufactured(CardPress, nil, []int{GoodPapyrus}, 1),

	// Civilian structures
	CardPawnshop: NewCardCivilian(CardPawnshop, nil, 3, nil, nil),
	CardBaths: NewCardCivilian(CardBaths, Cost{
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
	CardStockade: NewCardMilitary(CardStockade, Cost{
		GoodWood: 1,
	}, 1, nil, nil),
	CardBarracks: NewCardMilitary(CardBarracks, Cost{
		GoodOre: 1,
	}, 1, nil, nil),
	CardGuardTower: NewCardMilitary(CardGuardTower, Cost{
		GoodClay: 1,
	}, 1, nil, nil),

	// Scientific structures
	CardApothecary: NewCardScience(CardApothecary, Cost{
		GoodTextile: 1,
	}, FieldMathematics, nil, []string{
		CardStables,
		CardDispensary,
	}),
	CardWorkshop: NewCardScience(CardWorkshop, Cost{
		GoodGlass: 1,
	}, FieldEngineering, nil, []string{
		CardLaboratory,
		CardArcheryRange,
	}),
	CardScriptorium: NewCardScience(CardScriptorium, Cost{
		GoodPapyrus: 1,
	}, FieldTheology, nil, []string{
		CardCourthouse,
		CardLibrary,
	}),

	// Age 2

	// Raw goods
	CardSawmill: NewCardGoodRaw(CardSawmill, Cost{
		GoodCoin: 1,
	}, []int{GoodWood}, 2),
	CardQuarry: NewCardGoodRaw(CardQuarry, Cost{
		GoodCoin: 1,
	}, []int{GoodStone}, 2),
	CardBrickyard: NewCardGoodRaw(CardBrickyard, Cost{
		GoodCoin: 1,
	}, []int{GoodClay}, 2),
	CardFoundry: NewCardGoodRaw(CardFoundry, Cost{
		GoodCoin: 1,
	}, []int{GoodOre}, 2),

	// Civilian structures
	CardAqueduct: NewCardCivilian(CardAqueduct, Cost{
		GoodStone: 3,
	}, 5, []string{CardBaths}, nil),
	CardTemple: NewCardCivilian(CardTemple, Cost{
		GoodWood:  1,
		GoodClay:  1,
		GoodGlass: 1,
	}, 3, []string{CardAltar}, []string{CardPantheon}),
	CardStatue: NewCardCivilian(CardStatue, Cost{
		GoodOre:  2,
		GoodWood: 1,
	}, 4, []string{CardTheater}, []string{CardGardens}),
	CardCourthouse: NewCardCivilian(CardCourthouse, Cost{
		GoodClay:    2,
		GoodTextile: 1,
	}, 4, []string{CardScriptorium}, nil),

	// Commercial structures
	CardForum: NewCardGoodCommercial(CardForum, Cost{
		GoodClay: 2,
	}, ManufacturedGoods, 1, []string{
		CardEastTradingPost,
		CardWestTradingPost,
	}, []string{CardHaven}),
	CardCaravansery: NewCardGoodCommercial(CardCaravansery, Cost{
		GoodWood: 2,
	}, RawGoods, 1, []string{CardMarketplace}, []string{CardLighthouse}),
	CardVineyard: NewCardBonus(CardVineyard, CardKindCommercial, nil, []int{
		CardKindRaw,
	}, DirAll, 0, 1, nil, nil),
	CardBazar: NewCardBonus(CardBazar, CardKindCommercial, nil, []int{
		CardKindManufactured,
	}, DirAll, 0, 2, nil, nil),

	// Military structures
	CardWalls: NewCardMilitary(CardWalls, Cost{
		GoodStone: 3,
	}, 2, nil, []string{CardFortifications}),
	CardTrainingGround: NewCardMilitary(CardTrainingGround, Cost{
		GoodOre:  2,
		GoodWood: 1,
	}, 2, nil, []string{CardCircus}),
	CardStables: NewCardMilitary(CardStables, Cost{
		GoodClay: 1,
		GoodWood: 1,
		GoodOre:  1,
	}, 2, []string{CardApothecary}, nil),
	CardArcheryRange: NewCardMilitary(CardArcheryRange, Cost{
		GoodWood: 2,
		GoodOre:  1,
	}, 2, []string{CardWorkshop}, nil),

	// Scientific structures
	CardDispensary: NewCardScience(CardDispensary, Cost{
		GoodOre:   2,
		GoodGlass: 1,
	}, FieldMathematics, []string{CardApothecary}, []string{
		CardLodge,
		CardArena,
	}),
	CardLaboratory: NewCardScience(CardLaboratory, Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
	}, FieldEngineering, []string{CardWorkshop}, []string{
		CardObservatory,
		CardSiegeWorkshop,
	}),
	CardLibrary: NewCardScience(CardLibrary, Cost{
		GoodStone:   2,
		GoodTextile: 1,
	}, FieldTheology, []string{CardScriptorium}, []string{
		CardSenate,
		CardUniversity,
	}),
	CardSchool: NewCardScience(CardSchool, Cost{
		GoodWood:    1,
		GoodPapyrus: 1,
	}, FieldTheology, nil, []string{
		CardAcademy,
		CardStudy,
	}),

	// Age 3

	// Civilian structures
	CardPantheon: NewCardCivilian(CardPantheon, Cost{
		GoodClay:    2,
		GoodOre:     1,
		GoodGlass:   1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, 7, []string{CardTemple}, nil),
	CardGardens: NewCardCivilian(CardGardens, Cost{
		GoodClay: 2,
		GoodWood: 1,
	}, 5, []string{CardStatue}, nil),
	CardTownHall: NewCardCivilian(CardTownHall, Cost{
		GoodStone: 2,
		GoodOre:   1,
		GoodGlass: 1,
	}, 6, nil, nil),
	CardPalace: NewCardCivilian(CardPalace, Cost{
		GoodStone:   1,
		GoodOre:     1,
		GoodWood:    1,
		GoodClay:    1,
		GoodGlass:   1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, 8, nil, nil),
	CardSenate: NewCardCivilian(CardSenate, Cost{
		GoodWood:  2,
		GoodStone: 1,
		GoodOre:   1,
	}, 6, []string{CardLibrary}, nil),

	// Commercial structures
	CardHaven: NewCardBonus(CardHaven, CardKindCommercial, Cost{
		GoodWood:    1,
		GoodOre:     1,
		GoodTextile: 1,
	}, []int{CardKindRaw}, DirSelf, 1, 1, []string{CardForum}, nil),
	CardLighthouse: NewCardBonus(CardLighthouse, CardKindCommercial, Cost{
		GoodStone: 1,
		GoodGlass: 1,
	}, []int{CardKindCommercial}, DirSelf, 1, 1, []string{
		CardCaravansery,
	}, nil),
	CardChamberOfCommerce: NewCardBonus(CardChamberOfCommerce, CardKindCommercial, Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
	}, []int{CardKindManufactured}, DirSelf, 2, 2, nil, nil),
	CardArena: NewCardBonus(CardArena, CardKindCommercial, Cost{
		GoodStone: 2,
		GoodOre:   1,
	}, []int{WonderStage}, DirSelf, 1, 3, []string{
		CardDispensary,
	}, nil),

	// Military structures
	CardFortifications: NewCardMilitary(CardFortifications, Cost{
		GoodOre:   3,
		GoodStone: 1,
	}, 3, []string{CardWalls}, nil),
	CardCircus: NewCardMilitary(CardCircus, Cost{
		GoodStone: 3,
		GoodOre:   1,
	}, 3, []string{CardTrainingGround}, nil),
	CardArsenal: NewCardMilitary(CardArsenal, Cost{
		GoodWood:    2,
		GoodOre:     1,
		GoodTextile: 1,
	}, 3, nil, nil),
	CardSiegeWorkshop: NewCardMilitary(CardSiegeWorkshop, Cost{
		GoodClay: 3,
		GoodWood: 1,
	}, 3, []string{CardLaboratory}, nil),

	// Scientific structures
	CardLodge: NewCardScience(CardLodge, Cost{
		GoodClay:    2,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, FieldMathematics, []string{CardDispensary}, nil),
	CardObservatory: NewCardScience(CardObservatory, Cost{
		GoodOre:     2,
		GoodGlass:   1,
		GoodTextile: 1,
	}, FieldEngineering, []string{CardLaboratory}, nil),
	CardUniversity: NewCardScience(CardUniversity, Cost{
		GoodWood:    2,
		GoodPapyrus: 1,
		GoodGlass:   1,
	}, FieldTheology, []string{CardLibrary}, nil),
	CardAcademy: NewCardScience(CardAcademy, Cost{
		GoodStone: 3,
		GoodGlass: 1,
	}, FieldMathematics, []string{CardSchool}, nil),
	CardStudy: NewCardScience(CardStudy, Cost{
		GoodWood:    1,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, FieldEngineering, []string{CardSchool}, nil),

	// Guilds

	CardWorkersGuild: NewCardBonus(CardWorkersGuild, CardKindGuild, Cost{
		GoodOre:   2,
		GoodClay:  1,
		GoodStone: 1,
		GoodWood:  1,
	}, []int{CardKindRaw}, DirNeighbours, 1, 0, nil, nil),
	CardCraftsmensGuild: NewCardBonus(CardCraftsmensGuild, CardKindGuild, Cost{
		GoodOre:   2,
		GoodStone: 2,
	}, []int{CardKindManufactured}, DirNeighbours, 2, 0, nil, nil),
	CardTradersGuild: NewCardBonus(CardTradersGuild, CardKindGuild, Cost{
		GoodGlass:   1,
		GoodTextile: 1,
		GoodPapyrus: 1,
	}, []int{CardKindCommercial}, DirNeighbours, 1, 0, nil, nil),
	CardPhilosophersGuild: NewCardBonus(CardPhilosophersGuild, CardKindGuild, Cost{
		GoodClay:    3,
		GoodPapyrus: 1,
		GoodTextile: 1,
	}, []int{CardKindScientific}, DirNeighbours, 1, 0, nil, nil),
	CardSpiesGuild: NewCardBonus(CardSpiesGuild, CardKindGuild, Cost{
		GoodClay:  3,
		GoodGlass: 1,
	}, []int{CardKindMilitary}, DirNeighbours, 1, 0, nil, nil),
	CardStrategistsGuild: NewCardBonus(CardStrategistsGuild, CardKindGuild, Cost{
		GoodOre:     2,
		GoodStone:   1,
		GoodTextile: 1,
	}, []int{TokenDefeat}, DirNeighbours, 1, 0, nil, nil),
	CardShipownersGuild: NewCardBonus(CardShipownersGuild, CardKindGuild, Cost{
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
			Cost{
				GoodWood:    2,
				GoodOre:     2,
				GoodPapyrus: 1,
			},
			nil,
			nil,
		},
		AllFields,
	},
	CardMagistratesGuild: NewCardBonus(CardMagistratesGuild, CardKindGuild, Cost{
		GoodWood:    3,
		GoodStone:   1,
		GoodTextile: 1,
	}, []int{CardKindCivilian}, DirNeighbours, 1, 0, nil, nil),
	CardBuildersGuild: NewCardBonus(CardBuildersGuild, CardKindGuild, Cost{
		GoodStone: 2,
		GoodClay:  2,
		GoodGlass: 1,
	}, []int{WonderStage}, DirAll, 1, 0, nil, nil),
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
