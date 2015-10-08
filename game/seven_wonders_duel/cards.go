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
	CardSchool
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

var Cards map[int]Card

func init() {
	Cards = map[int]Card{
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
		CardSawmill: {
			Id:   CardSawmill,
			Name: "Sawmill",
			Type: CardTypeRaw,
			Cost: cost.Cost{GoodCoin: 2},
			Provides: []cost.Cost{
				{GoodWood: 2},
			},
		},
		CardBrickyard: {
			Id:   CardBrickyard,
			Name: "Brickyard",
			Type: CardTypeRaw,
			Cost: cost.Cost{GoodCoin: 2},
			Provides: []cost.Cost{
				{GoodClay: 2},
			},
		},
		CardShelfQuarry: {
			Id:   CardShelfQuarry,
			Name: "Shelf Quarry",
			Type: CardTypeRaw,
			Cost: cost.Cost{GoodCoin: 2},
			Provides: []cost.Cost{
				{GoodStone: 2},
			},
		},
		CardGlassblower: {
			Id:   CardGlassblower,
			Name: "Glassblower",
			Type: CardTypeManufactured,
			Provides: []cost.Cost{
				{GoodGlass: 1},
			},
		},
		CardDryingRoom: {
			Id:   CardDryingRoom,
			Name: "Drying Room",
			Type: CardTypeManufactured,
			Provides: []cost.Cost{
				{GoodPapyrus: 1},
			},
		},
		CardWalls: {
			Id:       CardWalls,
			Name:     "Walls",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodStone: 2},
			Military: 2,
		},
		CardForum: {
			Id:   CardForum,
			Name: "Forum",
			Type: CardTypeCommercial,
			Cost: cost.Cost{GoodCoin: 3, GoodClay: 1},
			Provides: []cost.Cost{
				{GoodGlass: 1},
				{GoodPapyrus: 1},
			},
		},
		CardCaravansery: {
			Id:   CardCaravansery,
			Name: "Caravansery",
			Type: CardTypeCommercial,
			Cost: cost.Cost{GoodCoin: 2, GoodGlass: 1, GoodPapyrus: 1},
			Provides: []cost.Cost{
				{GoodWood: 1},
				{GoodClay: 1},
				{GoodStone: 1},
			},
		},
		CardCustomsHouse: {
			Id:       CardCustomsHouse,
			Name:     "Customs House",
			Type:     CardTypeCommercial,
			Cost:     cost.Cost{GoodCoin: 4},
			Cheapens: []int{GoodPapyrus, GoodGlass},
		},
		CardTribunal: {
			Id:    CardTribunal,
			Name:  "Tribunal",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodWood: 2, GoodGlass: 1},
			VPRaw: 5,
		},
		CardHorseBreeders: {
			Id:       CardHorseBreeders,
			Name:     "Horse Breeders",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodClay: 1, GoodWood: 1},
			Military: 1,
		},
		CardBarracks: {
			Id:       CardBarracks,
			Name:     "Barracks",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodCoin: 3},
			Military: 1,
		},
		CardArcheryRange: {
			Id:        CardArcheryRange,
			Name:      "Archery Range",
			Type:      CardTypeMilitary,
			Cost:      cost.Cost{GoodStone: 1, GoodWood: 1, GoodPapyrus: 1},
			Military:  2,
			MakesFree: CardSiegeWorkshop,
		},
		CardParadeGround: {
			Id:        CardParadeGround,
			Name:      "Parade Ground",
			Type:      CardTypeMilitary,
			Cost:      cost.Cost{GoodClay: 2, GoodGlass: 1},
			Military:  2,
			MakesFree: CardCircus,
		},
		CardLibrary: {
			Id:      CardLibrary,
			Name:    "Library",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodStone: 1, GoodWood: 1, GoodGlass: 1},
			Science: ScienceLiterature,
			VPRaw:   2,
		},
		CardDispensary: {
			Id:      CardDispensary,
			Name:    "Dispensary",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodClay: 2, GoodStone: 1},
			Science: ScienceMedicine,
			VPRaw:   2,
		},
		CardSchool: {
			Id:        CardSchool,
			Name:      "School",
			Type:      CardTypeScientific,
			Cost:      cost.Cost{GoodWood: 1, GoodPapyrus: 2},
			Science:   ScienceEngineering,
			VPRaw:     1,
			MakesFree: CardUniversity,
		},
		CardLaboratory: {
			Id:        CardLaboratory,
			Name:      "Labratory",
			Type:      CardTypeScientific,
			Cost:      cost.Cost{GoodWood: 1, GoodGlass: 2},
			Science:   ScienceEngineering,
			VPRaw:     1,
			MakesFree: CardObservatory,
		},
		CardStatue: {
			Id:        CardStatue,
			Name:      "Statue",
			Type:      CardTypeCivilian,
			Cost:      cost.Cost{GoodClay: 2},
			VPRaw:     4,
			MakesFree: CardGardens,
		},
		CardTemple: {
			Id:        CardTemple,
			Name:      "Temple",
			Type:      CardTypeCivilian,
			Cost:      cost.Cost{GoodWood: 1, GoodPapyrus: 1},
			VPRaw:     4,
			MakesFree: CardPantheon,
		},
		CardAqueduct: {
			Id:    CardAqueduct,
			Name:  "Aqueduct",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodStone: 3},
			VPRaw: 5,
		},
		CardRostrum: {
			Id:        CardRostrum,
			Name:      "Rostrum",
			Type:      CardTypeCivilian,
			Cost:      cost.Cost{GoodStone: 1, GoodWood: 1},
			VPRaw:     4,
			MakesFree: CardSenate,
		},
		CardBrewery: {
			Id:   CardBrewery,
			Name: "Brewery",
			Type: CardTypeCommercial,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(player, 6)
			},
			MakesFree: CardArena,
		},
		CardArsenal: {
			Id:       CardArsenal,
			Name:     "Arsenal",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodClay: 3, GoodWood: 2},
			Military: 3,
		},
		CardCourthouse: {
			Id:       CardCourthouse,
			Name:     "Courthouse",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodCoin: 8},
			Military: 3,
		},
		CardAcademy: {
			Id:      CardAcademy,
			Name:    "Academy",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodStone: 1, GoodWood: 1, GoodGlass: 2},
			Science: ScienceAstronomy,
			VPRaw:   3,
		},
		CardStudy: {
			Id:      CardStudy,
			Name:    "Study",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodWood: 2, GoodGlass: 1, GoodPapyrus: 1},
			Science: ScienceAstronomy,
			VPRaw:   3,
		},
		CardChamberOfCommerce: {
			Id:    CardChamberOfCommerce,
			Name:  "Chamber of Commerce",
			Type:  CardTypeCommercial,
			Cost:  cost.Cost{GoodPapyrus: 2},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					g.PlayerCardTypeCount(player, CardTypeManufactured)*3,
				)
			},
		},
		CardPort: {
			Id:    CardPort,
			Name:  "Port",
			Type:  CardTypeCommercial,
			Cost:  cost.Cost{GoodWood: 1, GoodGlass: 1, GoodPapyrus: 1},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					g.PlayerCardTypeCount(player, CardTypeRaw)*2,
				)
			},
		},
		CardArmory: {
			Id:    CardArmory,
			Name:  "Armory",
			Type:  CardTypeCommercial,
			Cost:  cost.Cost{GoodStone: 2, GoodGlass: 1},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					g.PlayerCardTypeCount(player, CardTypeMilitary),
				)
			},
		},
		CardPalace: {
			Id:   CardPalace,
			Name: "Palace",
			Type: CardTypeCivilian,
			Cost: cost.Cost{
				GoodClay:  1,
				GoodStone: 1,
				GoodWood:  1,
				GoodGlass: 2,
			},
			VPRaw: 7,
		},
		CardTownHall: {
			Id:    CardTownHall,
			Name:  "Town Hall",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodStone: 3, GoodWood: 2},
			VPRaw: 7,
		},
		CardObelisk: {
			Id:    CardObelisk,
			Name:  "Obelisk",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodStone: 2, GoodGlass: 1},
			VPRaw: 5,
		},
		CardFortifications: {
			Id:       CardFortifications,
			Name:     "Fortifications",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodStone: 2, GoodClay: 1, GoodPapyrus: 1},
			Military: 2,
		},
		CardSiegeWorkshop: {
			Id:       CardSiegeWorkshop,
			Name:     "Siege Workshop",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodWood: 3, GoodGlass: 1},
			Military: 2,
		},
		CardCircus: {
			Id:       CardCircus,
			Name:     "Circus",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodClay: 2, GoodStone: 2},
			Military: 2,
		},
		CardUniversity: {
			Id:      CardUniversity,
			Name:    "University",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodClay: 1, GoodGlass: 1, GoodPapyrus: 1},
			Science: ScienceCartography,
			VPRaw:   2,
		},
		CardObservatory: {
			Id:      CardObservatory,
			Name:    "Observatory",
			Type:    CardTypeScientific,
			Cost:    cost.Cost{GoodStone: 1, GoodPapyrus: 2},
			Science: ScienceCartography,
			VPRaw:   2,
		},
		CardGardens: {
			Id:    CardGardens,
			Name:  "Gardens",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodClay: 2, GoodWood: 2},
			VPRaw: 6,
		},
		CardPantheon: {
			Id:    CardPantheon,
			Name:  "Pantheon",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodClay: 1, GoodWood: 1, GoodPapyrus: 2},
			VPRaw: 6,
		},
		CardSenate: {
			Id:    CardSenate,
			Name:  "Senate",
			Type:  CardTypeCivilian,
			Cost:  cost.Cost{GoodClay: 2, GoodStone: 1, GoodPapyrus: 1},
			VPRaw: 5,
		},
		CardLighthouse: {
			Id:    CardLighthouse,
			Name:  "Lighthouse",
			Type:  CardTypeCommercial,
			Cost:  cost.Cost{GoodClay: 2, GoodGlass: 1},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					g.PlayerCardTypeCount(player, CardTypeCommercial),
				)
			},
		},
		CardArena: {
			Id:    CardArena,
			Name:  "Arena",
			Type:  CardTypeCommercial,
			Cost:  cost.Cost{GoodClay: 1, GoodStone: 1, GoodWood: 1},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(player, len(g.PlayerWonders[player])*2)
			},
		},
	}
}
