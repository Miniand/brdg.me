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

func DeckAge2() card.Deck {
	d := card.Deck{}
	for _, c := range []card.Deck{
		// Raw goods
		NewCardGoodRaw(CardSawmill, Cost{
			GoodCoin: 1,
		}, []int{GoodWood}, 2, 3, 4),
		NewCardGoodRaw(CardQuarry, Cost{
			GoodCoin: 1,
		}, []int{GoodStone}, 2, 3, 4),
		NewCardGoodRaw(CardBrickyard, Cost{
			GoodCoin: 1,
		}, []int{GoodClay}, 2, 3, 4),
		NewCardGoodRaw(CardFoundry, Cost{
			GoodCoin: 1,
		}, []int{GoodOre}, 2, 3, 4),

		// Manufactured goods
		NewCardGoodManufactured(CardLoom, nil, []int{GoodTextile}, 1, 3, 5),
		NewCardGoodManufactured(CardGlassworks, nil, []int{GoodGlass}, 1, 3, 5),
		NewCardGoodManufactured(CardPress, nil, []int{GoodPapyrus}, 1, 3, 5),

		// Civilian structures
		NewCardCivilian(CardAqueduct, Cost{
			GoodStone: 3,
		}, 5, []string{CardBaths}, nil, 3, 7),
		NewCardCivilian(CardTemple, Cost{
			GoodWood:  1,
			GoodClay:  1,
			GoodGlass: 1,
		}, 3, []string{CardAltar}, []string{CardPantheon}, 3, 6),
		NewCardCivilian(CardStatue, Cost{
			GoodOre:  2,
			GoodWood: 1,
		}, 4, []string{CardTheater}, []string{CardGardens}, 3, 7),
		NewCardCivilian(CardCourthouse, Cost{
			GoodClay:    2,
			GoodTextile: 1,
		}, 4, []string{CardScriptorium}, nil, 3, 5),

		// Commercial structures
		NewCardGoodCommercial(CardForum, Cost{
			GoodClay: 2,
		}, ManufacturedGoods, 1, []string{
			CardEastTradingPost,
			CardWestTradingPost,
		}, []string{CardHaven}, 3, 6, 7),
		NewCardGoodCommercial(CardCaravansery, Cost{
			GoodWood: 2,
		}, RawGoods, 1, []string{CardMarketplace}, []string{CardLighthouse}, 3, 5, 6),
		NewCardBonus(CardVineyard, CardKindCommercial, nil, []int{
			CardKindRaw,
		}, DirAll, 0, 1, nil, nil, 3, 6),
		NewCardBonus(CardBazar, CardKindCommercial, nil, []int{
			CardKindManufactured,
		}, DirAll, 0, 2, nil, nil, 4, 7),

		// Military structures
		NewCardMilitary(CardWalls, Cost{
			GoodStone: 3,
		}, 2, nil, []string{CardFortifications}, 3, 7),
		NewCardMilitary(CardTrainingGround, Cost{
			GoodOre:  2,
			GoodWood: 1,
		}, 2, nil, []string{CardCircus}, 4, 6, 7),
		NewCardMilitary(CardStables, Cost{
			GoodClay: 1,
			GoodWood: 1,
			GoodOre:  1,
		}, 2, []string{CardApothecary}, nil, 3, 5),
		NewCardMilitary(CardArcheryRange, Cost{
			GoodWood: 2,
			GoodOre:  1,
		}, 2, []string{CardWorkshop}, nil, 3, 6),

		// Scientific structures
		NewCardScience(CardDispensary, Cost{
			GoodOre:   2,
			GoodGlass: 1,
		}, FieldMathematics, []string{CardApothecary}, []string{
			CardLodge,
			CardArena,
		}, 3, 4),
		NewCardScience(CardLaboratory, Cost{
			GoodClay:    2,
			GoodPapyrus: 1,
		}, FieldEngineering, []string{CardWorkshop}, []string{
			CardObservatory,
			CardSiegeWorkshop,
		}, 3, 5),
		NewCardScience(CardLibrary, Cost{
			GoodStone:   2,
			GoodTextile: 1,
		}, FieldTheology, []string{CardScriptorium}, []string{
			CardSenate,
			CardUniversity,
		}, 3, 6),
		NewCardScience(CardSchool, Cost{
			GoodWood:    1,
			GoodPapyrus: 1,
		}, FieldTheology, nil, []string{
			CardAcademy,
			CardStudy,
		}, 3, 7),
	} {
		d = d.PushMany(c)
	}
	return d
}

func DeckAge3() card.Deck {
	d := card.Deck{}
	for _, c := range []card.Deck{
		// Civilian structures
		NewCardCivilian(CardPantheon, Cost{
			GoodClay:    2,
			GoodOre:     1,
			GoodGlass:   1,
			GoodPapyrus: 1,
			GoodTextile: 1,
		}, 7, []string{CardTemple}, nil, 3, 6),
		NewCardCivilian(CardGardens, Cost{
			GoodClay: 2,
			GoodWood: 1,
		}, 5, []string{CardStatue}, nil, 3, 4),
		NewCardCivilian(CardTownHall, Cost{
			GoodStone: 2,
			GoodOre:   1,
			GoodGlass: 1,
		}, 6, nil, nil, 3, 5, 6),
		NewCardCivilian(CardPalace, Cost{
			GoodStone:   1,
			GoodOre:     1,
			GoodWood:    1,
			GoodClay:    1,
			GoodGlass:   1,
			GoodPapyrus: 1,
			GoodTextile: 1,
		}, 8, nil, nil, 3, 7),
		NewCardCivilian(CardSenate, Cost{
			GoodWood:  2,
			GoodStone: 1,
			GoodOre:   1,
		}, 6, []string{CardLibrary}, nil, 3, 5),

		// Commercial structures
		NewCardBonus(CardHaven, CardKindCommercial, Cost{
			GoodWood:    1,
			GoodOre:     1,
			GoodTextile: 1,
		}, []int{CardKindRaw}, DirSelf, 1, 1, []string{CardForum}, nil, 3, 4),
		NewCardBonus(CardLighthouse, CardKindCommercial, Cost{
			GoodStone: 1,
			GoodGlass: 1,
		}, []int{CardKindCommercial}, DirSelf, 1, 1, []string{
			CardCaravansery,
		}, nil, 3, 6),
		NewCardBonus(CardChamberOfCommerce, CardKindCommercial, Cost{
			GoodClay:    2,
			GoodPapyrus: 1,
		}, []int{CardKindManufactured}, DirSelf, 2, 2, nil, nil, 4, 6),
		NewCardBonus(CardArena, CardKindCommercial, Cost{
			GoodStone: 2,
			GoodOre:   1,
		}, []int{WonderStage}, DirSelf, 1, 3, []string{
			CardDispensary,
		}, nil, 3, 5, 7),

		// Military structures
		NewCardMilitary(CardFortifications, Cost{
			GoodOre:   3,
			GoodStone: 1,
		}, 3, []string{CardWalls}, nil, 3, 7),
		NewCardMilitary(CardCircus, Cost{
			GoodStone: 3,
			GoodOre:   1,
		}, 3, []string{CardTrainingGround}, nil, 4, 5, 6),
		NewCardMilitary(CardArsenal, Cost{
			GoodWood:    2,
			GoodOre:     1,
			GoodTextile: 1,
		}, 3, nil, nil, 3, 4, 7),
		NewCardMilitary(CardSiegeWorkshop, Cost{
			GoodClay: 3,
			GoodWood: 1,
		}, 3, []string{CardLaboratory}, nil, 3, 5),

		// Scientific structures
		NewCardScience(CardLodge, Cost{
			GoodClay:    2,
			GoodPapyrus: 1,
			GoodTextile: 1,
		}, FieldMathematics, []string{CardDispensary}, nil, 3, 6),
		NewCardScience(CardObservatory, Cost{
			GoodOre:     2,
			GoodGlass:   1,
			GoodTextile: 1,
		}, FieldEngineering, []string{CardLaboratory}, nil, 3, 7),
		NewCardScience(CardUniversity, Cost{
			GoodWood:    2,
			GoodPapyrus: 1,
			GoodGlass:   1,
		}, FieldTheology, []string{CardLibrary}, nil, 3, 4),
		NewCardScience(CardAcademy, Cost{
			GoodStone: 3,
			GoodGlass: 1,
		}, FieldMathematics, []string{CardSchool}, nil, 3, 7),
		NewCardScience(CardStudy, Cost{
			GoodWood:    1,
			GoodPapyrus: 1,
			GoodTextile: 1,
		}, FieldEngineering, []string{CardSchool}, nil, 3, 5),
	} {
		d = d.PushMany(c)
	}
	return d
}

func DeckGuild() card.Deck {
	return card.Deck{
		NewCardBonus(CardWorkersGuild, CardKindGuild, Cost{
			GoodOre:   2,
			GoodClay:  1,
			GoodStone: 1,
			GoodWood:  1,
		}, []int{CardKindRaw}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardCraftsmensGuild, CardKindGuild, Cost{
			GoodOre:   2,
			GoodStone: 2,
		}, []int{CardKindManufactured}, DirNeighbours, 2, 0, nil, nil, 3)[0],
		NewCardBonus(CardTradersGuild, CardKindGuild, Cost{
			GoodGlass:   1,
			GoodTextile: 1,
			GoodPapyrus: 1,
		}, []int{CardKindCommercial}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardPhilosophersGuild, CardKindGuild, Cost{
			GoodClay:    3,
			GoodPapyrus: 1,
			GoodTextile: 1,
		}, []int{CardKindScientific}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardSpiesGuild, CardKindGuild, Cost{
			GoodClay:  3,
			GoodGlass: 1,
		}, []int{CardKindMilitary}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardStrategistsGuild, CardKindGuild, Cost{
			GoodOre:     2,
			GoodStone:   1,
			GoodTextile: 1,
		}, []int{TokenDefeat}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardShipownersGuild, CardKindGuild, Cost{
			GoodWood:    3,
			GoodGlass:   1,
			GoodPapyrus: 1,
		}, []int{
			CardKindRaw,
			CardKindManufactured,
			CardKindGuild,
		}, DirSelf, 1, 0, nil, nil, 3)[0],
		CardScience{
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
				3,
			},
			AllFields,
		},
		NewCardBonus(CardMagistratesGuild, CardKindGuild, Cost{
			GoodWood:    3,
			GoodStone:   1,
			GoodTextile: 1,
		}, []int{CardKindCivilian}, DirNeighbours, 1, 0, nil, nil, 3)[0],
		NewCardBonus(CardBuildersGuild, CardKindGuild, Cost{
			GoodStone: 2,
			GoodClay:  2,
			GoodGlass: 1,
		}, []int{WonderStage}, DirAll, 1, 0, nil, nil, 3)[0],
	}
}
