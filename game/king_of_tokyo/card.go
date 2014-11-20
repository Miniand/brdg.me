package king_of_tokyo

const (
	CardKindKeep = iota
	CardKindDiscard
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

func Deck() []CardBase {
	return []CardBase{
		CardAcidAttack{},
		CardAlienMetabolism{},
		CardAlphaMonster{},
		CardApartmentBuilding{},
		CardArmorPlating{},
		CardBackgroundDweller{},
		CardBurrowing{},
		CardCamouflage{},
		CardCommuterTrain{},
		CardCompleteDestruction{},
		CardCornerStore{},
		CardDeathFromAbove{},
		CardDedicatedNewsTeam{},
		CardEaterOfTheDead{},
		CardEnergize{},
		CardEnergyHoarder{},
		CardEvacuationOrders{},
		CardEvacuationOrders{},
		CardEvenBigger{},
		CardExtraHead{},
		CardExtraHead{},
		CardFireBlast{},
		CardFireBreathing{},
		CardFreezeTime{},
		CardFrenzy{},
		CardFriendOfChildren{},
		CardGasRefinery{},
		CardGiantBrain{},
		CardGourmet{},
		CardHeal{},
		CardHealingRay{},
		CardHerbivore{},
		CardHerdCuller{},
		CardHighAltitudeBombing{},
		CardItHasAChild{},
		CardJetFighters{},
		CardJets{},
		CardMadeInALab{},
		CardMetamorph{},
		CardMimic{},
		CardMonsterBatteries{},
		CardNationalGuard{},
		CardNovaBreath{},
		CardNuclearPowerPlant{},
		CardOmnivore{},
		CardOpportunist{},
		CardParasiticTentacles{},
		CardPlotTwist{},
		CardPoisonQuills{},
		CardPoisonSpit{},
		CardPsychicProbe{},
		CardRapidHealing{},
		CardRegeneration{},
		CardRootingForTheUnderdog{},
		CardShrinkRay{},
		CardSkyscraper{},
		CardSmokeCloud{},
		CardSolarPowered{},
		CardSpikedTail{},
		CardStretchy{},
		CardTanks{},
		CardTelepath{},
		CardUrbavore{},
		CardVastStorm{},
		CardWereOnlyMakingItStronger{},
		CardWings{},
	}
}

func Shuffle(deck []CardBase) []CardBase {
	l := len(deck)
	shuffled := make([]CardBase, l)
	for i, p := range r.Perm(l) {
		shuffled[i] = deck[p]
	}
	return shuffled
}

type AttackModifier interface {
	ModifyAttack(game *Game, attack int) int
}

type HasThings interface {
	Things() []interface{}
}

type CardCostModifier interface {
	ModifyCardCost(game *Game, cost int) int
}

type PostCardBuyHandler interface {
	PostCardBuy(game *Game, card CardBase, cost int)
}

type PostAttackHandler interface {
	PostAttack(game *Game, attack int)
}
