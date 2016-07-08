package lords_of_vegas

type BoardSpace struct {
	Location      string
	PayCasino     int
	StartingMoney int
	Dice          int
	BuildPrice    int
	Strip         bool
}

type BoardSpaceState struct {
	Owned  bool
	Owner  int
	Dice   int
	Casino int
}

var BoardSpaces = []BoardSpace{
	{"A1", CasinoPioneer, 7, 3, 9, false},
	{"A2", CasinoAlbion, 8, 2, 6, false},
	{"A3", CasinoVega, 5, 5, 15, true},
	{"A4", CasinoSphynx, 6, 4, 12, false},
	{"A5", CasinoTivoli, 7, 3, 9, false},
	{"A6", CasinoTheStrip, 4, 6, 20, true},

	{"B1", CasinoSphynx, 5, 5, 15, true},
	{"B2", CasinoTivoli, 8, 2, 6, false},
	{"B3", CasinoAlbion, 7, 3, 9, false},
	{"B4", CasinoAlbion, 4, 6, 20, true},
	{"B5", CasinoVega, 7, 3, 9, false},
	{"B6", CasinoPioneer, 6, 4, 12, false},

	{"C1", CasinoTivoli, 6, 4, 12, false},
	{"C2", CasinoSphynx, 7, 3, 9, false},
	{"C3", CasinoAlbion, 4, 6, 20, true},
	{"C4", CasinoPioneer, 8, 2, 6, false},
	{"C5", CasinoVega, 9, 1, 8, false},
	{"C6", CasinoTivoli, 6, 4, 12, true},
	{"C7", CasinoSphynx, 8, 2, 6, false},
	{"C8", CasinoAlbion, 9, 1, 8, false},
	{"C9", CasinoPioneer, 6, 4, 12, true},
	{"C10", CasinoVega, 7, 3, 9, false},
	{"C11", CasinoTivoli, 8, 2, 6, false},
	{"C12", CasinoSphynx, 5, 5, 15, true},

	{"D1", CasinoTivoli, 4, 6, 20, true},
	{"D2", CasinoPioneer, 7, 3, 9, false},
	{"D3", CasinoVega, 6, 4, 12, false},
	{"D4", CasinoVega, 6, 4, 12, true},
	{"D5", CasinoTheStrip, 9, 1, 8, false},
	{"D6", CasinoAlbion, 8, 2, 6, false},
	{"D7", CasinoAlbion, 5, 5, 15, true},
	{"D8", CasinoSphynx, 8, 2, 6, false},
	{"D9", CasinoPioneer, 7, 3, 9, false},

	{"E1", CasinoPioneer, 7, 3, 9, false},
	{"E2", CasinoAlbion, 8, 2, 6, false},
	{"E3", CasinoVega, 5, 5, 15, true},
	{"E4", CasinoTivoli, 6, 4, 12, false},
	{"E5", CasinoPioneer, 7, 3, 9, false},
	{"E6", CasinoSphynx, 4, 6, 20, true},

	{"F1", CasinoAlbion, 4, 6, 20, true},
	{"F2", CasinoTivoli, 7, 3, 9, false},
	{"F3", CasinoSphynx, 6, 4, 12, false},
	{"F4", CasinoSphynx, 6, 4, 12, true},
	{"F5", CasinoPioneer, 9, 1, 8, false},
	{"F6", CasinoVega, 8, 2, 6, false},
	{"F7", CasinoVega, 5, 5, 15, true},
	{"F8", CasinoTheStrip, 8, 2, 6, false},
	{"F9", CasinoTivoli, 7, 3, 9, false},
}

var BoardLayout = [][]string{
	{"A1", "A2", "A3", "ST", "B1", "B2", "B3"},
	{"A4", "A5", "A6", "ST", "B4", "B5", "B6"},
	{"  ", "  ", "  ", "ST", "  ", "  ", "  "},
	{"C1", "C2", "C3", "ST", "D1", "D2", "D3"},
	{"C4", "C5", "C6", "ST", "D4", "D5", "D6"},
	{"C7", "C8", "C9", "ST", "D7", "D8", "D9"},
	{"C10", "C11", "C12", "ST", "  ", "  ", "  "},
	{"  ", "  ", "  ", "ST", "F1", "F2", "F3"},
	{"E1", "E2", "E3", "ST", "F4", "F5", "F6"},
	{"E4", "E5", "E6", "ST", "F7", "F8", "F9"},
}

var BoardSpaceByLocation = map[string]BoardSpace{}

func init() {
	for _, s := range BoardSpaces {
		BoardSpaceByLocation[s.Location] = s
	}
}
