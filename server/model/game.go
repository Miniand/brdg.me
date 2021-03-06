package model

import (
	"errors"
	"math/rand"
	"time"

	"github.com/Miniand/brdg.me/game"
	r "github.com/dancannon/gorethink"
)

func GameTable() r.Term {
	return r.Table("games")
}

func LoadGame(id string) (*GameModel, error) {
	m := &GameModel{
		Id: id,
	}
	return m, m.Load()
}

func GamesForPlayer(player string) (*r.Cursor, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return GameTable().GetAllByIndex("PlayerList", player).Run(session)
}

func ActiveGamesForPlayer(player string) (*r.Cursor, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return GameTable().GetAllByIndex(
		"IsFinished:PlayerList",
		[]interface{}{false, player},
	).Run(session)
}

func CurrentTurnGamesForPlayer(player string) (*r.Cursor, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return GameTable().GetAllByIndex(
		"IsFinished:WhoseTurn",
		[]interface{}{false, player},
	).OrderBy(r.Row.Field("WhoseTurnSince").Field(player)).Run(session)
}

func NotCurrentTurnGamesForPlayer(player string) (*r.Cursor, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return GameTable().GetAllByIndex(
		"IsFinished:PlayerList",
		[]interface{}{false, player},
	).Filter(func(row r.Term) interface{} {
		return r.Not(row.Field("WhoseTurn").Contains(player))
	}).Run(session)
}

func RecentlyFinishedGamesForPlayer(player string) (*r.Cursor, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return GameTable().OrderBy(r.OrderByOpts{Index: r.Desc("FinishedAt")}).
		Filter(r.Row.Field("IsFinished").Eq(true).And(
		r.Row.Field("PlayerList").Contains(player))).Limit(5).Run(session)
}

func SaveGame(g game.Playable) (*GameModel, error) {
	gm, err := GameToGameModel(g)
	if err != nil {
		return nil, err
	}
	err = gm.Save()
	return gm, err
}

func UpdateGame(id string, g game.Playable) (*GameModel, error) {
	gm, err := GameToGameModel(g)
	if err != nil {
		return nil, err
	}
	gm.Id = id
	err = gm.Save()
	return gm, err
}

func StartNewGame(g game.Playable, players []string) (*GameModel, error) {
	// Unique players
	playerMap := map[string]bool{}
	for _, p := range players {
		playerMap[p] = true
	}
	uniquePlayers := []string{}
	for p, _ := range playerMap {
		uniquePlayers = append(uniquePlayers, p)
	}
	// Shuffle players
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	l := len(uniquePlayers)
	perm := r.Perm(l)
	shuffledPlayers := make([]string, l)
	for i := 0; i < l; i++ {
		shuffledPlayers[i] = uniquePlayers[perm[i]]
	}
	// Start game
	if err := g.Start(shuffledPlayers); err != nil {
		return nil, err
	}
	// Save game
	return SaveGame(g)
}

func GameToGameModel(g game.Playable) (*GameModel, error) {
	gm := &GameModel{}
	if err := gm.UpdateState(g); err != nil {
		return nil, err
	}
	return gm, nil
}

type GameModel struct {
	Id                   string `gorethink:"id,omitempty"`
	PlayerList           []string
	Winners              []string
	EliminatedPlayerList []string
	IsFinished           bool
	FinishedAt           time.Time `gorethink:",omitempty"`
	WhoseTurnSince       map[string]time.Time
	WhoseTurn            []string
	Type                 string
	State                []byte
	Restarted            bool
	ConcedePlayers       []string
	ConcedeVote          map[string]bool
}

func (gm *GameModel) ToGame() (game.Playable, error) {
	g := game.RawCollection()[gm.Type]
	if g == nil {
		return nil, errors.New("Unable to find game type " + gm.Type)
	}
	err := g.Decode(gm.State)
	return g, err
}

func (gm *GameModel) Finish(winners []string) {
	gm.IsFinished = true
	gm.FinishedAt = time.Now()
	gm.WhoseTurn = []string{}
	gm.WhoseTurnSince = map[string]time.Time{}
	gm.Winners = winners
}

func (gm *GameModel) UpdateState(g game.Playable) error {
	state, err := g.Encode()
	if err != nil {
		return err
	}
	gm.State = state
	gm.Type = g.Identifier()
	gm.PlayerList = g.PlayerList()
	isFinished := g.IsFinished()
	if !gm.IsFinished && isFinished {
		gm.Finish(g.Winners())
	} else if !isFinished {
		if e, ok := g.(game.Eliminator); ok {
			gm.EliminatedPlayerList = e.EliminatedPlayerList()
		}
		// Cache whose turn it is and set the time for people whose turn it has
		// just become.
		if gm.IsConcedeVoting() {
			gm.WhoseTurn = gm.RemainingConcedeVotePlayers()
		} else {
			gm.WhoseTurn = g.WhoseTurn()
		}
		if gm.WhoseTurnSince == nil {
			gm.WhoseTurnSince = map[string]time.Time{}
		}
		whoseTurnMap := map[string]bool{}
		for _, p := range gm.WhoseTurn {
			whoseTurnMap[p] = true
			if gm.WhoseTurnSince[p].IsZero() {
				gm.WhoseTurnSince[p] = time.Now()
			}
		}
		for p, _ := range gm.WhoseTurnSince {
			if !whoseTurnMap[p] {
				delete(gm.WhoseTurnSince, p)
			}
		}
	}
	return nil
}

func (gm *GameModel) IsConcedeVoting() bool {
	return !gm.IsFinished && gm.ConcedeVote != nil
}

func (gm *GameModel) RemainingConcedeVotePlayers() []string {
	remaining := []string{}
	if !gm.IsConcedeVoting() {
		return remaining
	}
	eliminated := map[string]bool{}
	if gm.EliminatedPlayerList != nil {
		for _, ep := range gm.EliminatedPlayerList {
			eliminated[ep] = true
		}
	}
	for _, p := range gm.PlayerList {
		if _, ok := gm.ConcedeVote[p]; !ok && !eliminated[p] {
			remaining = append(remaining, p)
		}
	}
	return remaining
}

func (gm *GameModel) Load() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	res, err := GameTable().Get(gm.Id).Run(session)
	if err != nil {
		return err
	}
	return res.One(gm)
}

func (gm *GameModel) Save() error {
	var rqlTerm r.Term
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	if gm.Id == "" {
		rqlTerm = GameTable().Insert(gm)
	} else {
		rqlTerm = GameTable().Get(gm.Id).Update(gm)
	}
	res, err := rqlTerm.RunWrite(session)
	if err != nil {
		return err
	}
	if gm.Id == "" {
		gm.Id = res.GeneratedKeys[0]
	}
	return nil
}
