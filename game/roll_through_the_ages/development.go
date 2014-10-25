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
	DevelopmentCoinage
	DevelopmentCaravans
	DevelopmentReligion
	DevelopmentGranaries
	DevelopmentMasonry
	DevelopmentEngineering
	DevelopmentArchitecture
	DevelopmentEmpire
)

var Developments = []int{
	DevelopmentLeadership,
	DevelopmentIrrigation,
	DevelopmentAgriculture,
	DevelopmentQuarrying,
	DevelopmentMedicine,
	DevelopmentCoinage,
	DevelopmentCaravans,
	DevelopmentReligion,
	DevelopmentGranaries,
	DevelopmentMasonry,
	DevelopmentEngineering,
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
		Cost:   15,
		Points: 3,
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
	DevelopmentReligion: {
		Name:   "religion",
		Effect: "revolt affects opponents",
		Cost:   20,
		Points: 6,
	},
	DevelopmentGranaries: {
		Name:   "granaries",
		Effect: "sell food for 4 coins each",
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
	DevelopmentArchitecture: {
		Name:   "architecture",
		Effect: "bonus pts: 1 / monument",
		Cost:   50,
		Points: 8,
	},
	DevelopmentEmpire: {
		Name:   "empire",
		Effect: "bonus pts: 1 / city",
		Cost:   60,
		Points: 8,
	},
}
