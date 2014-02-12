package model

import (
	"errors"
	"github.com/Miniand/brdg.me/game"
	r "github.com/dancannon/gorethink"
	"math/rand"
	"time"
)

type GameModel struct {
	Id         string `gorethink:"id,omitempty"`
	PlayerList []string
	Winners    []string
	IsFinished bool
	WhoseTurn  []string
	Type       string
	State      []byte
}

func GameTable() r.RqlTerm {
	return r.Table("games")
}

func LoadGame(id string) (*GameModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	row, err := GameTable().Get(id).RunRow(session)
	if err != nil {
		return nil, err
	}
	m := &GameModel{}
	if err := row.Scan(m); err != nil {
		return nil, err
	}
	return m, nil
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
	state, err := g.Encode()
	if err != nil {
		return nil, err
	}
	gm := &GameModel{
		PlayerList: g.PlayerList(),
		Winners:    g.Winners(),
		IsFinished: g.IsFinished(),
		WhoseTurn:  g.WhoseTurn(),
		Type:       g.Identifier(),
		State:      state,
	}
	return gm, nil
}

func (gm *GameModel) ToGame() (game.Playable, error) {
	g := game.RawCollection()[gm.Type]
	if g == nil {
		return nil, errors.New("Unable to find game type " + gm.Type)
	}
	err := g.Decode(gm.State)
	return g, err
}

func (gm *GameModel) Save() error {
	var rqlTerm r.RqlTerm
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
