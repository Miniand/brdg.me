package seven_wonders_duel

import (
	"fmt"
	"strings"

	"github.com/Miniand/brdg.me/game/cost"
	"github.com/Miniand/brdg.me/game/helper"
)

const (
	CardTypeRaw = iota + 1
	CardTypeManufactured
	CardTypeCivilian
	CardTypeScientific
	CardTypeCommercial
	CardTypeMilitary
	CardTypeGuild
	CardTypeWonder
	CardTypeProgress
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

	// Wonders
	WonderTheAppianWay
	WonderTheMausoleum
	WonderCircusMaximus
	WonderPiraeus
	WonderTheColossus
	WonderThePyramids
	WonderTheGreatLibrary
	WonderTheSphinx
	WonderTheGreatLighthouse
	WonderTheStatueOfZeus
	WonderTheHangingGardens
	WonderTheTempleOfArtemis

	// Progress tokens
	ProgressAgriculture
	ProgressArchitecture
	ProgressEconomy
	ProgressLaw
	ProgressMasonry
	ProgressMathematics
	ProgressPhilosophy
	ProgressStrategy
	ProgressTheology
	ProgressUrbanism
)

type Card struct {
	Id            int
	Name          string
	Type          int
	Cost          cost.Cost
	VPRaw         int
	VPFunc        func(g *Game, player int) int
	Coin          int
	AfterBuild    func(g *Game, player int)
	Provides      []cost.Cost
	MakesFree     int
	Military      int
	Science       int
	Cheapens      []int
	DiscountGoods func(player, card int) int
	ExtraTurn     bool
	Summary       string
}

func Age1Cards() []int {
	return []int{
		CardLumberYard,
		CardLoggingCamp,
		CardClayPool,
		CardClayPit,
		CardQuarry,
		CardStonePit,
		CardGlassworks,
		CardPress,
		CardGuardTower,
		CardWorkshop,
		CardApothecary,
		CardStoneReserve,
		CardClayReserve,
		CardWoodReserve,
		CardStable,
		CardGarrison,
		CardPalisade,
		CardScriptorium,
		CardPharmacist,
		CardTheater,
		CardAltar,
		CardBaths,
		CardTavern,
	}
}

func Age2Cards() []int {
	return []int{
		CardSawmill,
		CardBrickyard,
		CardShelfQuarry,
		CardGlassblower,
		CardDryingRoom,
		CardWalls,
		CardForum,
		CardCaravansery,
		CardCustomsHouse,
		CardTribunal,
		CardHorseBreeders,
		CardBarracks,
		CardArcheryRange,
		CardParadeGround,
		CardLibrary,
		CardDispensary,
		CardSchool,
		CardLaboratory,
		CardStatue,
		CardTemple,
		CardAqueduct,
		CardRostrum,
		CardBrewery,
	}
}

func Age3Cards() []int {
	return []int{
		CardArsenal,
		CardCourthouse,
		CardAcademy,
		CardStudy,
		CardChamberOfCommerce,
		CardPort,
		CardArmory,
		CardPalace,
		CardTownHall,
		CardObelisk,
		CardFortifications,
		CardSiegeWorkshop,
		CardCircus,
		CardUniversity,
		CardObservatory,
		CardGardens,
		CardPantheon,
		CardSenate,
		CardLighthouse,
		CardArena,
	}
}

func Wonders() []int {
	return []int{
		WonderTheAppianWay,
		WonderCircusMaximus,
		WonderTheColossus,
		WonderTheGreatLibrary,
		WonderTheGreatLighthouse,
		WonderTheHangingGardens,
		WonderTheMausoleum,
		WonderPiraeus,
		WonderThePyramids,
		WonderTheSphinx,
		WonderTheStatueOfZeus,
		WonderTheTempleOfArtemis,
	}
}

func ProgressTokens() []int {
	return []int{
		ProgressAgriculture,
		ProgressArchitecture,
		ProgressEconomy,
		ProgressLaw,
		ProgressMasonry,
		ProgressMathematics,
		ProgressPhilosophy,
		ProgressStrategy,
		ProgressTheology,
		ProgressUrbanism,
	}
}

type Layout [][]int

type Loc struct {
	X, Y int
}

func (l Layout) IsVisible(loc Loc) bool {
	if loc.Y%2 == 0 {
		return true
	}
	return l.CanBuild(loc)
}

func (l Layout) CanBuild(loc Loc) bool {
	for _, bl := range l.LocsBelow(loc) {
		if l.At(bl) != 0 {
			return false
		}
	}
	return true
}

func (l Layout) LocsBelow(loc Loc) []Loc {
	offset := 0
	if loc.Y%2 == 0 {
		offset = -1
	}
	return []Loc{
		{loc.X + offset, loc.Y + 1},
		{loc.X + 1 + offset, loc.Y + 1},
	}
}

func (l Layout) At(loc Loc) int {
	if loc.X < 0 || loc.Y < 0 || loc.Y >= len(l) || loc.X >= len(l[loc.Y]) {
		return 0
	}
	return l[loc.Y][loc.X]
}

var AgeLayouts = map[int]func() Layout{
	1: Age1Layout,
	2: Age2Layout,
	3: Age3Layout,
}

func Age1Layout() Layout {
	cards := helper.IntShuffle(Age1Cards())
	layout := Layout{}
	for i := 0; i < 5; i++ {
		row := make([]int, (4-i)/2)
		n := 2 + i
		row = append(row, cards[:n]...)
		cards = cards[n:]
		layout = append(layout, row)
	}
	return layout
}

func Age2Layout() Layout {
	cards := helper.IntShuffle(Age2Cards())
	layout := Layout{}
	for i := 0; i < 5; i++ {
		row := make([]int, (i / 2))
		n := 6 - i
		row = append(row, cards[:n]...)
		cards = cards[n:]
		layout = append(layout, row)
	}
	return layout
}

func Age3Layout() Layout {
	cards := helper.IntShuffle(Age3Cards())[:17]
	cards = append(cards, helper.IntShuffle(GuildCards())[:3]...)
	cards = helper.IntShuffle(cards)
	layout := Layout{}
	for i := 0; i < 3; i++ {
		row := make([]int, (2-i)/2)
		n := 2 + i
		row = append(row, cards[:n]...)
		cards = cards[n:]
		layout = append(layout, row)
	}
	layout = append(layout, []int{cards[0], 0, cards[1]})
	cards = cards[2:]
	for i := 0; i < 3; i++ {
		row := make([]int, (i / 2))
		n := 4 - i
		row = append(row, cards[:n]...)
		cards = cards[n:]
		layout = append(layout, row)
	}
	return layout
}

func GuildCards() []int {
	return []int{
		CardMerchantsGuild,
		CardShipownersGuild,
		CardBuildersGuild,
		CardMagistratesGuild,
		CardScientistsGuild,
		CardMoneylendersGuild,
		CardTacticiansGuild,
	}
}

func (c Card) VP(g *Game, player int) int {
	vp := c.VPRaw
	if c.VPFunc != nil {
		vp += c.VPFunc(g, player)
	}
	return vp
}

func (c Card) RenderSummary() string {
	if c.Summary != "" {
		return c.Summary
	}
	parts := []string{}
	if c.Provides != nil {
		parts = append(parts, RenderProvides(c.Provides))
	}
	if c.Military > 0 {
		parts = append(parts, RenderMilitary(c.Military))
	}
	if c.Science > 0 {
		parts = append(parts, RenderScience(c.Science))
	}
	if len(c.Cheapens) > 0 {
		parts = append(parts, fmt.Sprintf(
			"%s costs %s",
			strings.Join(RenderGoods(c.Cheapens), " "),
			RenderCoins(1),
		))
	}
	if c.Coin > 0 {
		parts = append(parts, RenderCoins(c.Coin))
	}
	if c.VPRaw > 0 {
		parts = append(parts, RenderVP(c.VPRaw))
	}
	if c.ExtraTurn {
		parts = append(parts, ExtraTurnText)
	}
	if c.MakesFree > 0 {
		parts = append(parts, LinkBuildText)
	}
	return strings.Join(parts, ", ")
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
			Name:     "Stone Res.",
			Type:     CardTypeCommercial,
			Cost:     cost.Cost{GoodCoin: 3},
			Cheapens: []int{GoodStone},
		},
		CardClayReserve: {
			Id:       CardClayReserve,
			Name:     "Clay Res.",
			Type:     CardTypeCommercial,
			Cost:     cost.Cost{GoodCoin: 3},
			Cheapens: []int{GoodClay},
		},
		CardWoodReserve: {
			Id:       CardWoodReserve,
			Name:     "Wood Res.",
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
			Id:        CardTavern,
			Name:      "Tavern",
			Type:      CardTypeCommercial,
			Coin:      4,
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
			Name:     "Customs",
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
			Name:     "Horse Breed",
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
			Name:      "Archery Rng.",
			Type:      CardTypeMilitary,
			Cost:      cost.Cost{GoodStone: 1, GoodWood: 1, GoodPapyrus: 1},
			Military:  2,
			MakesFree: CardSiegeWorkshop,
		},
		CardParadeGround: {
			Id:        CardParadeGround,
			Name:      "Parade Grnd.",
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
			Id:        CardBrewery,
			Name:      "Brewery",
			Type:      CardTypeCommercial,
			Coin:      6,
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
			Id:   CardChamberOfCommerce,
			Name: "Chamber",
			Type: CardTypeCommercial,
			Summary: fmt.Sprintf(
				"%s x %s, %s",
				RenderCoins(3),
				RenderCardType(CardTypeManufactured),
				RenderVP(3),
			),
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
			Id:   CardPort,
			Name: "Port",
			Type: CardTypeCommercial,
			Summary: fmt.Sprintf(
				"%s x %s, %s",
				RenderCoins(2),
				RenderCardType(CardTypeRaw),
				RenderVP(3),
			),
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
			Id:   CardArmory,
			Name: "Armory",
			Type: CardTypeCommercial,
			Summary: fmt.Sprintf(
				"%s x %s, %s",
				RenderCoins(1),
				RenderCardType(CardTypeMilitary),
				RenderVP(3),
			),
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
			Name:     "Fort",
			Type:     CardTypeMilitary,
			Cost:     cost.Cost{GoodStone: 2, GoodClay: 1, GoodPapyrus: 1},
			Military: 2,
		},
		CardSiegeWorkshop: {
			Id:       CardSiegeWorkshop,
			Name:     "Siege WkShop",
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
			Id:   CardLighthouse,
			Name: "Lighthouse",
			Type: CardTypeCommercial,
			Summary: fmt.Sprintf(
				"%s x %s, %s",
				RenderCoins(1),
				RenderCardType(CardTypeCommercial),
				RenderVP(3),
			),
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
			Id:   CardArena,
			Name: "Arena",
			Type: CardTypeCommercial,
			Summary: fmt.Sprintf(
				"%s x %s, %s",
				RenderCoins(2),
				RenderCardType(CardTypeWonder),
				RenderVP(3),
			),
			Cost:  cost.Cost{GoodClay: 1, GoodStone: 1, GoodWood: 1},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					g.PlayerCardTypeCount(player, CardTypeWonder)*2,
				)
			},
		},
		CardMerchantsGuild: {
			Id:   CardMerchantsGuild,
			Name: "Merchants",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s %s ^ %s",
				RenderVP(1),
				RenderCoins(1),
				RenderCardType(CardTypeCommercial),
			),
			Cost: cost.Cost{
				GoodClay:    1,
				GoodWood:    1,
				GoodGlass:   1,
				GoodPapyrus: 1,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(CardTypeCommercial)
			},
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					Cards[CardMerchantsGuild].VPFunc(g, player),
				)
			},
		},
		CardShipownersGuild: {
			Id:   CardShipownersGuild,
			Name: "Shipowners",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s %s ^ %s %s",
				RenderVP(1),
				RenderCoins(1),
				RenderCardType(CardTypeRaw),
				RenderCardType(CardTypeManufactured),
			),
			Cost: cost.Cost{
				GoodClay:    1,
				GoodStone:   1,
				GoodGlass:   1,
				GoodPapyrus: 1,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(
					CardTypeRaw,
					CardTypeManufactured,
				)
			},
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					Cards[CardShipownersGuild].VPFunc(g, player),
				)
			},
		},
		CardBuildersGuild: {
			Id:   CardBuildersGuild,
			Name: "Builders",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s ^ %s",
				RenderVP(2),
				WonderText,
			),
			Cost: cost.Cost{
				GoodStone: 2,
				GoodClay:  1,
				GoodWood:  1,
				GoodGlass: 1,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(CardTypeWonder) * 2
			},
		},
		CardMagistratesGuild: {
			Id:   CardMagistratesGuild,
			Name: "Magistrates",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s %s ^ %s",
				RenderVP(1),
				RenderCoins(1),
				RenderCardType(CardTypeCivilian),
			),
			Cost: cost.Cost{
				GoodWood:    2,
				GoodClay:    1,
				GoodPapyrus: 1,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(CardTypeCivilian)
			},
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					Cards[CardMagistratesGuild].VPFunc(g, player),
				)
			},
		},
		CardScientistsGuild: {
			Id:   CardScientistsGuild,
			Name: "Scientists",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s %s ^ %s",
				RenderVP(1),
				RenderCoins(1),
				RenderCardType(CardTypeScientific),
			),
			Cost: cost.Cost{
				GoodClay: 2,
				GoodWood: 2,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(CardTypeScientific)
			},
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					Cards[CardScientistsGuild].VPFunc(g, player),
				)
			},
		},
		CardMoneylendersGuild: {
			Id:   CardMoneylendersGuild,
			Name: "Moneylenders",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s ^ %s",
				RenderVP(1),
				RenderCoins(1),
			),
			Cost: cost.Cost{
				GoodStone: 2,
				GoodWood:  2,
			},
			VPFunc: func(g *Game, player int) int {
				return helper.IntMax(
					g.PlayerCoins[0],
					g.PlayerCoins[1],
				) / 3
			},
		},
		CardTacticiansGuild: {
			Id:   CardTacticiansGuild,
			Name: "Tacticians",
			Type: CardTypeGuild,
			Summary: fmt.Sprintf(
				"%s %s ^ %s",
				RenderVP(1),
				RenderCoins(1),
				RenderCardType(CardTypeMilitary),
			),
			Cost: cost.Cost{
				GoodStone:   2,
				GoodClay:    1,
				GoodPapyrus: 1,
			},
			VPFunc: func(g *Game, player int) int {
				return g.GreatestCardCount(CardTypeMilitary)
			},
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(
					player,
					Cards[CardTacticiansGuild].VPFunc(g, player),
				)
			},
		},
		WonderTheAppianWay: {
			Id:   WonderTheAppianWay,
			Name: "Appian Way",
			Type: CardTypeWonder,
			Summary: fmt.Sprintf(
				"%s, %s, opp. %s",
				RenderCoins(3),
				ExtraTurnText,
				RenderCoins(-3),
			),
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodClay:    2,
				GoodStone:   2,
			},
			VPRaw: 3,
			AfterBuild: func(g *Game, player int) {
				g.ModifyCoins(Opponent(player), -3)
			},
			ExtraTurn: true,
		},
		WonderTheMausoleum: {
			Id:      WonderTheMausoleum,
			Name:    "Mausoleum",
			Type:    CardTypeWonder,
			Summary: "build disc. card",
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodGlass:   2,
				GoodClay:    2,
			},
			VPRaw: 2,
			AfterBuild: func(g *Game, player int) {
				panic("implement free build of discarded card")
			},
		},
		WonderCircusMaximus: {
			Id:      WonderCircusMaximus,
			Name:    "Circus Maximus",
			Type:    CardTypeWonder,
			Summary: fmt.Sprintf("remove opp. %s", RenderCardType(CardTypeManufactured)),
			Cost: cost.Cost{
				GoodGlass: 1,
				GoodWood:  1,
				GoodStone: 2,
			},
			VPRaw:    3,
			Military: 1,
			AfterBuild: func(g *Game, player int) {
				panic("implement discarding of opponent manufacturing card")
			},
		},
		WonderPiraeus: {
			Id:   WonderPiraeus,
			Name: "Piraeus",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodClay:  1,
				GoodStone: 1,
				GoodWood:  2,
			},
			VPRaw: 2,
			Provides: []cost.Cost{
				{GoodPapyrus: 1},
				{GoodGlass: 1},
			},
			ExtraTurn: true,
		},
		WonderTheColossus: {
			Id:   WonderTheColossus,
			Name: "Colossus",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodGlass: 1,
				GoodClay:  3,
			},
			VPRaw:    3,
			Military: 2,
		},
		WonderThePyramids: {
			Id:   WonderThePyramids,
			Name: "Pyramids",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodStone:   3,
			},
			VPRaw: 9,
		},
		WonderTheGreatLibrary: {
			Id:      WonderTheGreatLibrary,
			Name:    "Great Library",
			Type:    CardTypeWonder,
			Summary: fmt.Sprintf("get disc. %s", ProgressTokenText),
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodGlass:   1,
				GoodWood:    3,
			},
			VPRaw: 4,
			AfterBuild: func(g *Game, player int) {
				panic("implement randomly picking 3 discarded progress tokens and choosing one")
			},
		},
		WonderTheSphinx: {
			Id:   WonderTheSphinx,
			Name: "Sphinx",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodGlass: 2,
				GoodClay:  1,
				GoodStone: 1,
			},
			VPRaw:     6,
			ExtraTurn: true,
		},
		WonderTheGreatLighthouse: {
			Id:   WonderTheGreatLighthouse,
			Name: "Great Lighthouse",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodPapyrus: 2,
				GoodStone:   1,
				GoodWood:    1,
			},
			VPRaw: 4,
			Provides: []cost.Cost{
				{GoodWood: 1},
				{GoodStone: 1},
				{GoodClay: 1},
			},
		},
		WonderTheStatueOfZeus: {
			Id:      WonderTheStatueOfZeus,
			Name:    "Statue of Zeus",
			Type:    CardTypeWonder,
			Summary: fmt.Sprintf("remove opp. %s", RenderCardType(CardTypeRaw)),
			Cost: cost.Cost{
				GoodPapyrus: 2,
				GoodClay:    1,
				GoodWood:    1,
				GoodStone:   1,
			},
			VPRaw:    3,
			Military: 1,
			AfterBuild: func(g *Game, player int) {
				panic("implement discarding of opponent's raw good card")
			},
		},
		WonderTheHangingGardens: {
			Id:   WonderTheHangingGardens,
			Name: "Hanging Gardens",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodGlass:   1,
				GoodWood:    2,
			},
			VPRaw:     3,
			Coin:      6,
			ExtraTurn: true,
		},
		WonderTheTempleOfArtemis: {
			Id:   WonderTheTempleOfArtemis,
			Name: "Temple of Artemis",
			Type: CardTypeWonder,
			Cost: cost.Cost{
				GoodPapyrus: 1,
				GoodGlass:   1,
				GoodStone:   1,
				GoodWood:    1,
			},
			Coin:      12,
			ExtraTurn: true,
		},
		ProgressAgriculture: {
			Id:    ProgressAgriculture,
			Name:  "Agriculture",
			Type:  CardTypeProgress,
			Coin:  6,
			VPRaw: 4,
		},
		ProgressArchitecture: {
			Id:      ProgressArchitecture,
			Name:    "Architecture",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("%s 2 fewer res.", WonderText),
			DiscountGoods: func(player, card int) int {
				if Cards[card].Type == CardTypeWonder {
					return 2
				}
				return 0
			},
		},
		ProgressEconomy: {
			Id:      ProgressEconomy,
			Name:    "Economy",
			Type:    CardTypeProgress,
			Summary: "Take trade money",
		},
		ProgressLaw: {
			Id:      ProgressLaw,
			Name:    "Law",
			Type:    CardTypeProgress,
			Science: ScienceLaw,
		},
		ProgressMasonry: {
			Id:      ProgressMasonry,
			Name:    "Masonry",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("%s 2 fewer res.", RenderCardType(CardTypeCivilian)),
			DiscountGoods: func(player, card int) int {
				if Cards[card].Type == CardTypeCivilian {
					return 2
				}
				return 0
			},
		},
		ProgressMathematics: {
			Id:      ProgressMathematics,
			Name:    "Mathematics",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("%s x 3%s", RenderVP(3), RenderCardType(CardTypeProgress)),
		},
		ProgressPhilosophy: {
			Id:    ProgressPhilosophy,
			Name:  "Philosophy",
			Type:  CardTypeProgress,
			VPRaw: 7,
		},
		ProgressStrategy: {
			Id:      ProgressStrategy,
			Name:    "Strategy",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("Each %s +%s", RenderCardType(CardTypeMilitary), RenderMilitary(1)),
		},
		ProgressTheology: {
			Id:      ProgressTheology,
			Name:    "Theology",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("All %s give %s", WonderText, ExtraTurnText),
		},
		ProgressUrbanism: {
			Id:      ProgressUrbanism,
			Name:    "Urbanism",
			Type:    CardTypeProgress,
			Summary: fmt.Sprintf("%s, %s x %s build", RenderCoins(6), RenderCoins(4), LinkBuildText),
		},
	}
}
