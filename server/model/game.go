package model

import (
	"errors"
	"github.com/beefsack/brdg.me/game"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type GameModel struct {
	Id         interface{} "_id"
	PlayerList []string
	Winners    []string
	IsFinished bool
	WhoseTurn  []string
	Type       string
	State      string
}

func GameCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DatabaseName()).C("games")
}

func LoadGame(id interface{}) (*GameModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	m := &GameModel{}
	err = GameCollection(session).FindId(id).One(m)
	if m.Id == nil {
		m = nil
	}
	return m, err
}

func SaveGame(g game.Playable) (*GameModel, error) {
	gm, err := GameToGameModel(g)
	if err != nil {
		return nil, err
	}
	err = gm.Save()
	return gm, err
}

func UpdateGame(id interface{}, g game.Playable) (*GameModel, error) {
	gm, err := GameToGameModel(g)
	if err != nil {
		return nil, err
	}
	gm.Id = id
	err = gm.Save()
	return gm, err
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
		State:      string(state),
	}
	return gm, nil
	return nil, nil
}

func (gm *GameModel) ToGame() (game.Playable, error) {
	g := game.RawCollection()[gm.Type]
	if g == nil {
		return nil, errors.New("Unable to find game type " + gm.Type)
	}
	err := g.Decode([]byte(gm.State))
	return g, err
}

func (gm *GameModel) Save() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	if gm.Id == nil {
		gm.Id = bson.NewObjectId()
	}
	_, err = GameCollection(session).UpsertId(gm.Id, gm)
	if err != nil {
		return err
	}
	return nil
}
