package bang_dice

type CharBase interface {
	Name() string
	Description() string
	StartingLife() int
}

var Chars = []CharBase{
	CharBartCassidy{},
	CharBlackJack{},
	CharCalamityJanet{},
	CharElGringo{},
	CharJesseJones{},
	CharJourdonnais{},
	CharKitCarlson{},
	CharLuckyDuke{},
	CharPaulRegret{},
	CharPedroRamirez{},
	CharRoseDoolan{},
	CharSidKetchum{},
	CharSlabTheKiller{},
	CharSuzyLafayette{},
	CharVultureSam{},
	CharWillyTheKid{},
}
