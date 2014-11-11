package roll_through_the_ages

type Development struct {
	Name   string
	Effect string
	Cost   int
	Points int
}

const (
	DevelopmentLeadership = iota
	DevelopmentIrrigation
	DevelopmentAgriculture
	DevelopmentQuarrying
	DevelopmentMedicine
	DevelopmentPreservation
	DevelopmentCoinage
	DevelopmentCaravans
	DevelopmentShipping
	DevelopmentSmithing
	DevelopmentReligion
	DevelopmentGranaries
	DevelopmentMasonry
	DevelopmentEngineering
	DevelopmentCommerce
	DevelopmentArchitecture
	DevelopmentEmpire
)

var Developments = []int{
	DevelopmentLeadership,
	DevelopmentIrrigation,
	DevelopmentAgriculture,
	DevelopmentQuarrying,
	DevelopmentMedicine,
	DevelopmentPreservation,
	DevelopmentCoinage,
	DevelopmentCaravans,
	DevelopmentShipping,
	DevelopmentSmithing,
	DevelopmentReligion,
	DevelopmentGranaries,
	DevelopmentMasonry,
	DevelopmentEngineering,
	DevelopmentCommerce,
	DevelopmentArchitecture,
	DevelopmentEmpire,
}

var DevelopmentValues = map[int]Development{
	DevelopmentLeadership: {
		Name:   "leadership",
		Effect: "reroll 1 die (after last roll)",
		Cost:   10,
		Points: 2,
	},
	DevelopmentIrrigation: {
		Name:   "irrigation",
		Effect: "drought has no effect",
		Cost:   10,
		Points: 2,
	},
	DevelopmentAgriculture: {
		Name:   "agriculture",
		Effect: "+1 food / food die",
		Cost:   15,
		Points: 3,
	},
	DevelopmentQuarrying: {
		Name:   "quarrying",
		Effect: "+1 stone if collecting stone",
		Cost:   15,
		Points: 3,
	},
	DevelopmentMedicine: {
		Name:   "medicine",
		Effect: "pestilence has no effect",
		Cost:   20,
		Points: 4,
	},
	DevelopmentPreservation: {
		Name:   "preservation",
		Effect: "food x2 before roll for 1 pottery",
		Cost:   20,
		Points: 4,
	},
	DevelopmentCoinage: {
		Name:   "coinage",
		Effect: "coin die results are worth 12",
		Cost:   20,
		Points: 4,
	},
	DevelopmentCaravans: {
		Name:   "caravans",
		Effect: "no need to discard goods",
		Cost:   20,
		Points: 4,
	},
	DevelopmentShipping: {
		Name:   "shipping",
		Effect: "swap 1 good / ship",
		Cost:   25,
		Points: 5,
	},
	DevelopmentSmithing: {
		Name:   "smithing",
		Effect: "invasion affects opponents",
		Cost:   25,
		Points: 5,
	},
	DevelopmentReligion: {
		Name:   "religion",
		Effect: "revolt affects opponents",
		Cost:   25,
		Points: 7,
	},
	DevelopmentGranaries: {
		Name:   "granaries",
		Effect: "sell food for 6 coins each",
		Cost:   30,
		Points: 6,
	},
	DevelopmentMasonry: {
		Name:   "masonry",
		Effect: "+1 worker / worker die",
		Cost:   30,
		Points: 6,
	},
	DevelopmentEngineering: {
		Name:   "engineering",
		Effect: "use stone for 3 workers each",
		Cost:   40,
		Points: 6,
	},
	DevelopmentCommerce: {
		Name:   "commerce",
		Effect: "bonus pts: 1 / good",
		Cost:   40,
		Points: 8,
	},
	DevelopmentArchitecture: {
		Name:   "architecture",
		Effect: "bonus pts: 2 / monument",
		Cost:   60,
		Points: 8,
	},
	DevelopmentEmpire: {
		Name:   "empire",
		Effect: "bonus pts: 1 / city",
		Cost:   70,
		Points: 10,
	},
}

func DevelopmentNameMap() map[int]string {
	m := map[int]string{}
	for _, d := range Developments {
		m[d] = DevelopmentValues[d].Name
	}
	return m
}
