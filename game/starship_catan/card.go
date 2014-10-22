package starship_catan

import (
	"github.com/Miniand/brdg.me/command"
	"github.com/Miniand/brdg.me/game/card"
)

type Commander interface {
	Commands() []command.Command
}

type FullStringer interface {
	FullString() string
}

type UnsortableCard struct{}

func (c UnsortableCard) Compare(other card.Comparer) (int, bool) {
	return 0, false
}

func StartingCards() card.Deck {
	return card.Deck{
		ColonyCard{
			Name:     "Alioth VIII",
			Resource: ResourceCarbon,
			Dice:     1,
		},
		ColonyCard{
			Name:     "Megrez VII",
			Resource: ResourceFuel,
			Dice:     1,
		},
	}
}

func SectorBaseCards() card.Deck {
	return card.Deck{
		// Adventure planets
		AdventurePlanetCard{
			Name: AdventurePlanetHades,
		},
		AdventurePlanetCard{
			Name: AdventurePlanetPallas,
		},
		AdventurePlanetCard{
			Name: AdventurePlanetPicasso,
		},
		AdventurePlanetCard{
			Name: AdventurePlanetPoseidon,
		},
		// Colony planets
		ColonyCard{
			Name:     "Dubhe IV",
			Resource: ResourceCarbon,
			Dice:     2,
		},
		ColonyCard{
			Name:     "Phekda VI",
			Resource: ResourceFood,
			Dice:     1,
		},
		ColonyCard{
			Name:     "Merak V",
			Resource: ResourceFood,
			Dice:     3,
		},
		ColonyCard{
			Name:     "Alkor III",
			Resource: ResourceFuel,
			Dice:     3,
		},
		ColonyCard{
			Name:     "Bellatrix I",
			Resource: ResourceOre,
			Dice:     1,
		},
		ColonyCard{
			Name:     "Heka II",
			Resource: ResourceOre,
			Dice:     2,
		},
		// Pirates
		PirateCard{
			Strength: 2,
			Ransom:   3,
		},
		PirateCard{
			Strength: 3,
			Ransom:   3,
		},
		// Trading post planets
		TradeCard{
			Name:        "Alnitak IX",
			Resources:   []int{ResourceTrade},
			Price:       3,
			TradingPost: true,
		},
		TradeCard{
			Name:        "Beteigeuze VI",
			Resources:   []int{ResourceCarbon},
			Price:       3,
			TradingPost: true,
		},
		TradeCard{
			Name:        "Aigel X",
			Resources:   []int{ResourceOre},
			Price:       3,
			TradingPost: true,
		},
		TradeCard{
			Name:        "Mintaka II",
			Resources:   []int{ResourceFuel},
			Price:       3,
			TradingPost: true,
		},
		TradeCard{
			Name:        "Saiph VI",
			Resources:   []int{ResourceFood},
			Price:       3,
			TradingPost: true,
		},
		// Trade planets
		TradeCard{
			Name:      "Corendium VII",
			Resources: []int{ResourceCarbon},
			Price:     1,
		},
		TradeCard{
			Name:      "Tostoku I",
			Resources: []int{ResourceCarbon},
			Price:     2,
		},
		TradeCard{
			Name:      "Marsitis VI",
			Resources: []int{ResourceCarbon},
			Price:     4,
		},
		TradeCard{
			Name:      "Quartzee X",
			Resources: []int{ResourceCarbon},
			Price:     5,
		},
		TradeCard{
			Name:      "Planctoinis VII",
			Resources: []int{ResourceFood},
			Price:     1,
		},
		TradeCard{
			Name:      "Sputsallia IV",
			Resources: []int{ResourceFood},
			Price:     2,
		},
		TradeCard{
			Name:      "Pobeckifiked VI",
			Resources: []int{ResourceFood},
			Price:     4,
		},
		TradeCard{
			Name:      "Califasperum V",
			Resources: []int{ResourceFood},
			Price:     5,
		},
		TradeCard{
			Name:      "Litigus IX",
			Resources: []int{ResourceFuel},
			Price:     1,
		},
		TradeCard{
			Name:      "Gonsarium II",
			Resources: []int{ResourceFuel},
			Price:     2,
		},
		TradeCard{
			Name:      "Brocollar II",
			Resources: []int{ResourceFuel},
			Price:     4,
		},
		TradeCard{
			Name:      "Phlatiarum V",
			Resources: []int{ResourceFuel},
			Price:     5,
		},
		TradeCard{
			Name:      "Ireoni VII",
			Resources: []int{ResourceOre},
			Price:     1,
		},
		TradeCard{
			Name:      "Cupperius IV",
			Resources: []int{ResourceOre},
			Price:     2,
		},
		TradeCard{
			Name:      "Leedsi X",
			Resources: []int{ResourceOre},
			Price:     4,
		},
		TradeCard{
			Name:      "Bazaltide IV",
			Resources: []int{ResourceOre},
			Price:     5,
		},
		TradeCard{
			Name:      "Martkwal VIII",
			Resources: []int{ResourceTrade},
			Price:     1,
		},
		TradeCard{
			Name:      "Beowulf's Bane",
			Resources: []int{ResourceTrade},
			Price:     2,
		},
		TradeCard{
			Name:      "Parapeckis VII",
			Resources: []int{ResourceTrade},
			Price:     4,
		},
		TradeCard{
			Name:      "Martiin - Tempest II",
			Resources: []int{ResourceTrade},
			Price:     5,
		},
		TradeCard{
			Name:      "Kopernikus II",
			Resources: []int{ResourceScience},
			Price:     3,
		},
		TradeCard{
			Name:      "Diplomat Outpost",
			Resources: Goods,
			Price:     3,
			Direction: TradeDirBuy,
			Maximum:   1,
		},
		TradeCard{
			Name:      "Diplomat Outpost",
			Resources: Goods,
			Price:     3,
			Direction: TradeDirBuy,
			Maximum:   1,
		},
	}
}

func Sector1Cards() card.Deck {
	return card.Deck{
		EmptyCard{},
		EmptyCard{},
		TradeCard{
			Name:        "Green Folk Outpost",
			Resources:   []int{ResourceScience},
			Price:       4,
			Direction:   TradeDirSell,
			Maximum:     1,
			TradingPost: true,
		},
		TradeCard{
			Name:        "Diplomat Outpost",
			Resources:   Goods,
			Price:       3,
			Direction:   TradeDirBuy,
			Maximum:     1,
			TradingPost: true,
		},
		PirateCard{
			Strength: 2,
			Ransom:   3,
		},
		PirateCard{
			Strength: 3,
			Ransom:   3,
		},
		PirateCard{
			Strength: 4,
			Ransom:   3,
		},
	}
}

func Sector2Cards() card.Deck {
	return card.Deck{
		EmptyCard{},
		EmptyCard{},
		ColonyCard{
			Name:     "Benet-Nash IX",
			Resource: ResourceCarbon,
			Dice:     3,
		},
		ColonyCard{
			Name:     "Mizar X",
			Resource: ResourceFood,
			Dice:     2,
		},
		TradeCard{
			Name:        "Scientist Outpost",
			Resources:   []int{ResourceScience},
			Price:       2,
			Direction:   TradeDirBuy,
			Maximum:     1,
			TradingPost: true,
		},
		PirateCard{
			Strength:      4,
			Ransom:        5,
			DestroyCannon: true,
		},
		PirateCard{
			Strength:      5,
			Ransom:        5,
			DestroyModule: true,
		},
	}
}

func Sector3Cards() card.Deck {
	return card.Deck{
		EmptyCard{},
		EmptyCard{},
		ColonyCard{
			Name:     "Enif I",
			Resource: ResourceFuel,
			Dice:     2,
		},
		ColonyCard{
			Name:     "Theta Pegasi II",
			Resource: ResourceOre,
			Dice:     3,
		},
		TradeCard{
			Name:      "Merchant Outpost",
			Resources: Goods,
			Price:     3,
			Direction: TradeDirSell,
			Maximum:   2,
		},
		PirateCard{
			Strength:      5,
			Ransom:        5,
			DestroyModule: true,
		},
		PirateCard{
			Strength:      6,
			Ransom:        5,
			DestroyModule: true,
		},
	}
}

func Sector4Cards() card.Deck {
	return card.Deck{
		EmptyCard{},
		EmptyCard{},
		EmptyCard{},
		EmptyCard{},
		MedianCard{},
		PirateCard{
			Strength:      6,
			Ransom:        5,
			DestroyModule: true,
		},
		PirateCard{
			Strength:      7,
			Ransom:        5,
			DestroyModule: true,
		},
	}
}

func ShuffledSectorCards() card.Deck {
	return Sector4Cards().Shuffle().
		PushMany(Sector3Cards().Shuffle()).
		PushMany(Sector2Cards().Shuffle()).
		PushMany(Sector1Cards().Shuffle()).
		PushMany(SectorBaseCards().Shuffle())
}

func Adventure1Cards() card.Deck {
	return card.Deck{
		AdventureEnvironmentalCrisis{},
		AdventureDiplomaticGift{},
		AdventureMerchantGift{},
	}
}

func Adventure2Cards() card.Deck {
	return card.Deck{
		AdventureFamine{},
		AdventureWholesaleOrder1{},
		AdventurePirateNest{},
	}
}

func Adventure3Cards() card.Deck {
	return card.Deck{
		AdventureCouncilMeeting{},
		AdventureEpidemic{},
		AdventureEmergency{},
	}
}

func Adventure4Cards() card.Deck {
	return card.Deck{
		AdventureReconstruction{},
		AdventureMonument{},
		AdventureWholesaleOrder2{},
	}
}

func ShuffledAdventureCards() card.Deck {
	return Adventure4Cards().Shuffle().
		PushMany(Adventure3Cards().Shuffle()).
		PushMany(Adventure2Cards().Shuffle()).
		PushMany(Adventure1Cards().Shuffle())
}
