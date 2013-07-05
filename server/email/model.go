package main

import (
	"errors"
	"github.com/beefsack/boredga.me/game"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
)

type GameModel struct {
	Id         interface{}
	PlayerList []string
	Winners    []string
	IsFinished bool
	WhoseTurn  []string
	Type       string
	State      string
}

func Connect() (*mgo.Session, error) {
	addr := os.Getenv("BOREDGAME_MONGODB_ADDRESS")
	if addr == "" {
		addr = "localhost"
	}
	return mgo.Dial(addr)
}

func Collection(session *mgo.Session) *mgo.Collection {
	db := os.Getenv("BOREDGAME_MONGODB_DATABASE")
	if db == "" {
		db = "boredgame"
	}
	return session.DB(db).C("games")
}

func LoadGame(id interface{}) (*GameModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	m := &GameModel{}
	err = Collection(session).FindId(id).One(m)
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
	if gm.Id == nil {
		gm.Id = bson.NewObjectId()
	}
	selector := bson.M{
		"_id": gm.Id,
	}
	_, err = Collection(session).Upsert(selector, gm)
	return err
}
