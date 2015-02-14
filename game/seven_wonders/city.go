package seven_wonders

import "github.com/Miniand/brdg.me/game/cost"

var Cities = map[string]City{
	CityRhodesA: {
		CityRhodesA,
		GoodOre,
		[]string{
			WonderStageRhodesA1,
			WonderStageRhodesA2,
			WonderStageRhodesA3,
		},
	},
	CityRhodesB: {
		CityRhodesB,
		GoodOre,
		[]string{
			WonderStageRhodesB1,
			WonderStageRhodesB2,
		},
	},

	CityAlexandriaA: {
		CityAlexandriaA,
		GoodGlass,
		[]string{
			WonderStageAlexandriaA1,
			WonderStageAlexandriaA2,
			WonderStageAlexandriaA3,
		},
	},
	CityAlexandriaB: {
		CityAlexandriaB,
		GoodGlass,
		[]string{
			WonderStageAlexandriaB1,
			WonderStageAlexandriaB2,
			WonderStageAlexandriaB3,
		},
	},

	CityEphesusA: {
		CityEphesusA,
		GoodPapyrus,
		[]string{
			WonderStageEphesusA1,
			WonderStageEphesusA2,
			WonderStageEphesusA3,
		},
	},
	CityEphesusB: {
		CityEphesusB,
		GoodPapyrus,
		[]string{
			WonderStageEphesusB1,
			WonderStageEphesusB2,
			WonderStageEphesusB3,
		},
	},

	CityBabylonA: {
		CityBabylonA,
		GoodClay,
		[]string{
			WonderStageBabylonA1,
			WonderStageBabylonA2,
			WonderStageBabylonA3,
		},
	},
	CityBabylonB: {
		CityBabylonB,
		GoodClay,
		[]string{
			WonderStageBabylonB1,
			WonderStageBabylonB2,
			WonderStageBabylonB3,
		},
	},

	CityOlympiaA: {
		CityOlympiaA,
		GoodWood,
		[]string{
			WonderStageOlympiaA1,
			WonderStageOlympiaA2,
			WonderStageOlympiaA3,
		},
	},
	CityOlympiaB: {
		CityOlympiaB,
		GoodWood,
		[]string{
			WonderStageOlympiaB1,
			WonderStageOlympiaB2,
			WonderStageOlympiaB3,
		},
	},

	CityHalicarnassusA: {
		CityHalicarnassusA,
		GoodTextile,
		[]string{
			WonderStageHalicarnassusA1,
			WonderStageHalicarnassusA2,
			WonderStageHalicarnassusA3,
		},
	},
	CityHalicarnassusB: {
		CityHalicarnassusB,
		GoodTextile,
		[]string{
			WonderStageHalicarnassusB1,
			WonderStageHalicarnassusB2,
			WonderStageHalicarnassusB3,
		},
	},

	CityGizaA: {
		CityGizaA,
		GoodStone,
		[]string{
			WonderStageGizaA1,
			WonderStageGizaA2,
			WonderStageGizaA3,
		},
	},
	CityGizaB: {
		CityGizaB,
		GoodStone,
		[]string{
			WonderStageGizaB1,
			WonderStageGizaB2,
			WonderStageGizaB3,
			WonderStageGizaB4,
		},
	},
}

type City struct {
	Name            string
	InitialResource int
	WonderStages    []string
}

func (c City) GoodsProduced() []cost.Cost {
	return []cost.Cost{{c.InitialResource: 1}}
}

func (c City) GoodsTraded() []cost.Cost {
	return c.GoodsProduced()
}
