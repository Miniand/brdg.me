package seven_wonders_duel

import "github.com/Miniand/brdg.me/game/cost"

const (
	CardTypeRaw = iota + 1
	CardTypeManufactured
	CardTypeCivilian
	CardTypeScientific
	CardTypeCommercial
	CardTypeMilitary
	CardTypeGuild
)

const (
	GoodCoin = iota + 1
	GoodWood
	GoodClay
	GoodStone
	GoodGlass
	GoodPapyrus
)

const (
	ScienceCartography = iota + 1
	ScienceLaw
	ScienceAstronomy
	ScienceMathematics
	ScienceMedicine
	ScienceLiterature
	ScienceEngineering
)

const (
	// Age 1
	CardLumberYard = iota + 1
	CardLoggingCamp
	CardClayPool
	CardClayPit
	CardQuarry
	CardStonePit
	CardGlassworks
	CardPress
	CardGuardTower
	CardWorkshop
	CardApothecary
	CardStoneReserve
	CardClayReserve
	CardWoodReserve
	CardStable
	CardGarrison
	CardPalisade
	CardScriptorium
	CardPharmacist
	CardTheater
	CardAltar
	CardBaths
	CardTavern

	// Age 2
	CardSawmill
	CardBrickyard
	CardShelfQuarry
	CardGlassblower
	CardDryingRoom
	CardWalls
	CardForum
	CardCaravansery
	CardCustomsHouse
	CardTribunal
	CardHorseBreeders
	CardBarracks
	CardArcheryRange
	CardParadeGround
	CardLibrary
	CardDispensary
	CardLaboratory
	CardStatue
	CardTemple
	CardAqueduct
	CardRostrum
	CardBrewery

	// Age 3
	CardArsenal
	CardCourthouse
	CardAcademy
	CardStudy
	CardChamberOfCommerce
	CardPort
	CardArmory
	CardPalace
	CardTownHall
	CardObelisk
	CardFortifications
	CardSiegeWorkshop
	CardCircus
	CardUniversity
	CardObservatory
	CardGardens
	CardPantheon
	CardSenate
	CardLighthouse
	CardArena

	// Guild
	CardMerchantsGuild
	CardShipownersGuild
	CardBuildersGuild
	CardMagistratesGuild
	CardScientistsGuild
	CardMoneylendersGuild
	CardTacticiansGuild
)

type Card struct {
	Id         int
	Name       string
	Type       int
	Cost       cost.Cost
	VPRaw      int
	VPFunc     func(g *Game, player int) int
	AfterBuild func(g *Game, player int)
	Provides   []cost.Cost
	MakesFree  int
	Military   int
	Science    int
	Cheapens   []int
}

func (c Card) VP(g *Game, player int) int {
	vp := c.VPRaw
	if c.VPFunc != nil {
		vp += c.VPFunc(g, player)
	}
	return vp
}

var Cards = map[int]Card{
	CardLumberYard: {
		Id:   CardLumberYard,
		Name: "Lumber Yard",
		Type: CardTypeRaw,
		Provides: []cost.Cost{
			{GoodWood: 1},
		},
	},
	CardLoggingCamp: {
		Id:   CardLoggingCamp,
		Name: "Logging Camp",
		Type: CardTypeRaw,
		Cost: cost.Cost{GoodCoin: 1},
		Provides: []cost.Cost{
			{GoodWood: 1},
		},
	},
	CardClayPool: {
		Id:   CardClayPool,
		Name: "Clay Pool",
		Type: CardTypeRaw,
		Provides: []cost.Cost{
			{GoodClay: 1},
		},
	},
	CardClayPit: {
		Id:   CardClayPit,
		Name: "Clay Pit",
		Type: CardTypeRaw,
		Cost: cost.Cost{GoodCoin: 1},
		Provides: []cost.Cost{
			{GoodClay: 1},
		},
	},
	CardQuarry: {
		Id:   CardQuarry,
		Name: "Quarry",
		Type: CardTypeRaw,
		Provides: []cost.Cost{
			{GoodStone: 1},
		},
	},
	CardStonePit: {
		Id:   CardStonePit,
		Name: "Stone Pit",
		Type: CardTypeRaw,
		Cost: cost.Cost{GoodCoin: 1},
		Provides: []cost.Cost{
			{GoodStone: 1},
		},
	},
	CardGlassworks: {
		Id:   CardGlassworks,
		Name: "Glassworks",
		Type: CardTypeManufactured,
		Cost: cost.Cost{GoodCoin: 1},
		Provides: []cost.Cost{
			{GoodGlass: 1},
		},
	},
	CardPress: {
		Id:   CardPress,
		Name: "Press",
		Type: CardTypeManufactured,
		Cost: cost.Cost{GoodCoin: 1},
		Provides: []cost.Cost{
			{GoodPapyrus: 1},
		},
	},
	CardGuardTower: {
		Id:       CardGuardTower,
		Name:     "Guard Tower",
		Type:     CardTypeMilitary,
		Military: 1,
	},
	CardWorkshop: {
		Id:      CardWorkshop,
		Name:    "Workshop",
		Type:    CardTypeScientific,
		Cost:    cost.Cost{GoodPapyrus: 1},
		VPRaw:   1,
		Science: ScienceMathematics,
	},
	CardApothecary: {
		Id:      CardApothecary,
		Name:    "Apothecary",
		Type:    CardTypeScientific,
		Cost:    cost.Cost{GoodGlass: 1},
		VPRaw:   1,
		Science: ScienceEngineering,
	},
	CardStoneReserve: {
		Id:       CardStoneReserve,
		Name:     "Stone Reserve",
		Type:     CardTypeCommercial,
		Cost:     cost.Cost{GoodCoin: 3},
		Cheapens: []int{GoodStone},
	},
	CardClayReserve: {
		Id:       CardClayReserve,
		Name:     "Clay Reserve",
		Type:     CardTypeCommercial,
		Cost:     cost.Cost{GoodCoin: 3},
		Cheapens: []int{GoodClay},
	},
	CardWoodReserve: {
		Id:       CardWoodReserve,
		Name:     "Wood Reserve",
		Type:     CardTypeCommercial,
		Cost:     cost.Cost{GoodCoin: 3},
		Cheapens: []int{GoodWood},
	},
	CardStable: {
		Id:        CardStable,
		Name:      "Stable",
		Type:      CardTypeMilitary,
		Cost:      cost.Cost{GoodWood: 1},
		Military:  1,
		MakesFree: CardHorseBreeders,
	},
	CardGarrison: {
		Id:        CardGarrison,
		Name:      "Garrison",
		Type:      CardTypeMilitary,
		Cost:      cost.Cost{GoodClay: 1},
		Military:  1,
		MakesFree: CardBarracks,
	},
	CardPalisade: {
		Id:        CardPalisade,
		Name:      "Palisade",
		Type:      CardTypeMilitary,
		Cost:      cost.Cost{GoodCoin: 2},
		Military:  1,
		MakesFree: CardFortifications,
	},
	CardScriptorium: {
		Id:        CardScriptorium,
		Name:      "Scriptorium",
		Type:      CardTypeScientific,
		Cost:      cost.Cost{GoodCoin: 2},
		Science:   ScienceLiterature,
		MakesFree: CardLibrary,
	},
	CardPharmacist: {
		Id:        CardPharmacist,
		Name:      "Pharmacist",
		Type:      CardTypeScientific,
		Cost:      cost.Cost{GoodCoin: 2},
		Science:   ScienceMedicine,
		MakesFree: CardDispensary,
	},
	CardTheater: {
		Id:        CardTheater,
		Name:      "Theater",
		Type:      CardTypeCivilian,
		VPRaw:     3,
		MakesFree: CardStatue,
	},
	CardAltar: {
		Id:        CardAltar,
		Name:      "Altar",
		Type:      CardTypeCivilian,
		VPRaw:     3,
		MakesFree: CardTemple,
	},
	CardBaths: {
		Id:        CardBaths,
		Name:      "Baths",
		Type:      CardTypeCivilian,
		VPRaw:     3,
		MakesFree: CardAqueduct,
	},
	CardTavern: {
		Id:   CardTavern,
		Name: "Tavern",
		Type: CardTypeCommercial,
		AfterBuild: func(g *Game, player int) {
			g.ModifyCoins(player, 4)
		},
		MakesFree: CardLighthouse,
	},
}
