package king_of_tokyo

const (
	CardKindKeep = iota
	CardKindDiscard
)

const (
	AcidAttack = iota
	AlienMetabolism
	AlphaMonster
	ApartmentBuilding
	ArmorPlating
	BackgroundDweller
	Burrowing
	Camouflage
	CommuterTrain
	CompleteDestruction
	CornerStore
	DeathFromAbove
	DedicatedNewsTeam
	EaterOfTheDead
	Energize
	EnergyHoarder
	EvacuationOrders1
	EvacuationOrders2
	EvenBigger
	ExtraHead1
	ExtraHead2
	FireBlast
	FireBreathing
	FreezeTime
	Frenzy
	FriendOfChildren
	GasRefinery
	GiantBrain
	Gourmet
	Heal
	HealingRay
	Herbivore
	HerdCuller
	HighAltitudeBombing
	ItHasAChild
	JetFighters
	Jets
	MadeInALab
	Metamorph
	Mimic
	MonsterBatteries
	NationalGuard
	NovaBreath
	NuclearPowerPlant
	Omnivore
	Opportunist
	ParasiticTentacles
	PlotTwist
	PoisonQuills
	PoisonSpit
	PsychicProbe
	RapidHealing
	Regeneration
	RootingForTheUnderdog
	ShrinkRay
	Skyscraper
	SmokeCloud
	SolarPowered
	SpikedTail
	Stretchy
	Tanks
	Telepath
	Urbavore
	VastStorm
	WereOnlyMakingItStronger
	Wings
)

type Namer interface {
	Name() string
}

type Describer interface {
	Description() string
}

type Coster interface {
	Cost() int
}

type Kinder interface {
	Kind() int
}

type CardBase interface {
	Namer
	Describer
	Coster
	Kinder
}

var Cards = map[int]CardBase{
	AcidAttack:               &CardAcidAttack{},
	AlienMetabolism:          &CardAlienMetabolism{},
	AlphaMonster:             &CardAlphaMonster{},
	ApartmentBuilding:        &CardApartmentBuilding{},
	ArmorPlating:             &CardArmorPlating{},
	BackgroundDweller:        &CardBackgroundDweller{},
	Burrowing:                &CardBurrowing{},
	Camouflage:               &CardCamouflage{},
	CommuterTrain:            &CardCommuterTrain{},
	CompleteDestruction:      &CardCompleteDestruction{},
	CornerStore:              &CardCornerStore{},
	DeathFromAbove:           &CardDeathFromAbove{},
	DedicatedNewsTeam:        &CardDedicatedNewsTeam{},
	EaterOfTheDead:           &CardEaterOfTheDead{},
	Energize:                 &CardEnergize{},
	EnergyHoarder:            &CardEnergyHoarder{},
	EvacuationOrders1:        &CardEvacuationOrders{},
	EvacuationOrders2:        &CardEvacuationOrders{},
	EvenBigger:               &CardEvenBigger{},
	ExtraHead1:               &CardExtraHead{},
	ExtraHead2:               &CardExtraHead{},
	FireBlast:                &CardFireBlast{},
	FireBreathing:            &CardFireBreathing{},
	FreezeTime:               &CardFreezeTime{},
	Frenzy:                   &CardFrenzy{},
	FriendOfChildren:         &CardFriendOfChildren{},
	GasRefinery:              &CardGasRefinery{},
	GiantBrain:               &CardGiantBrain{},
	Gourmet:                  &CardGourmet{},
	Heal:                     &CardHeal{},
	HealingRay:               &CardHealingRay{},
	Herbivore:                &CardHerbivore{},
	HerdCuller:               &CardHerdCuller{},
	HighAltitudeBombing:      &CardHighAltitudeBombing{},
	ItHasAChild:              &CardItHasAChild{},
	JetFighters:              &CardJetFighters{},
	Jets:                     &CardJets{},
	MadeInALab:               &CardMadeInALab{},
	Metamorph:                &CardMetamorph{},
	Mimic:                    &CardMimic{},
	MonsterBatteries:         &CardMonsterBatteries{},
	NationalGuard:            &CardNationalGuard{},
	NovaBreath:               &CardNovaBreath{},
	NuclearPowerPlant:        &CardNuclearPowerPlant{},
	Omnivore:                 &CardOmnivore{},
	Opportunist:              &CardOpportunist{},
	ParasiticTentacles:       &CardParasiticTentacles{},
	PlotTwist:                &CardPlotTwist{},
	PoisonQuills:             &CardPoisonQuills{},
	PoisonSpit:               &CardPoisonSpit{},
	PsychicProbe:             &CardPsychicProbe{},
	RapidHealing:             &CardRapidHealing{},
	Regeneration:             &CardRegeneration{},
	RootingForTheUnderdog:    &CardRootingForTheUnderdog{},
	ShrinkRay:                &CardShrinkRay{},
	Skyscraper:               &CardSkyscraper{},
	SmokeCloud:               &CardSmokeCloud{},
	SolarPowered:             &CardSolarPowered{},
	SpikedTail:               &CardSpikedTail{},
	Stretchy:                 &CardStretchy{},
	Tanks:                    &CardTanks{},
	Telepath:                 &CardTelepath{},
	Urbavore:                 &CardUrbavore{},
	VastStorm:                &CardVastStorm{},
	WereOnlyMakingItStronger: &CardWereOnlyMakingItStronger{},
	Wings: &CardWings{},
}

var Deck = []int{
	AcidAttack,
	AlienMetabolism,
	AlphaMonster,
	ApartmentBuilding,
	ArmorPlating,
	BackgroundDweller,
	Burrowing,
	Camouflage,
	CommuterTrain,
	CompleteDestruction,
	CornerStore,
	DeathFromAbove,
	DedicatedNewsTeam,
	EaterOfTheDead,
	Energize,
	EnergyHoarder,
	EvacuationOrders1,
	EvacuationOrders2,
	EvenBigger,
	ExtraHead1,
	ExtraHead2,
	FireBlast,
	FireBreathing,
	FreezeTime,
	Frenzy,
	FriendOfChildren,
	GasRefinery,
	GiantBrain,
	Gourmet,
	Heal,
	HealingRay,
	Herbivore,
	HerdCuller,
	HighAltitudeBombing,
	ItHasAChild,
	JetFighters,
	Jets,
	MadeInALab,
	Metamorph,
	Mimic,
	MonsterBatteries,
	NationalGuard,
	NovaBreath,
	NuclearPowerPlant,
	Omnivore,
	Opportunist,
	ParasiticTentacles,
	PlotTwist,
	PoisonQuills,
	PoisonSpit,
	PsychicProbe,
	RapidHealing,
	Regeneration,
	RootingForTheUnderdog,
	ShrinkRay,
	Skyscraper,
	SmokeCloud,
	SolarPowered,
	SpikedTail,
	Stretchy,
	Tanks,
	Telepath,
	Urbavore,
	VastStorm,
	WereOnlyMakingItStronger,
	Wings,
}
