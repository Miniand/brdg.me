package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type UserModel struct {
	Id           interface{} "_id"
	Email        string
	Unsubscribed bool
}

func UserCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DatabaseName()).C("users")
}

func LoadUser(id interface{}) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	m := &UserModel{}
	err = UserCollection(session).FindId(id).One(m)
	if m.Id == nil {
		m = nil
	}
	return m, err
}

func FirstUserByEmail(email string) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	m := &UserModel{}
	err = UserCollection(session).Find(bson.M{
		"email": email,
	}).One(m)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	return m, err
}

func (um *UserModel) Save() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	if um.Id == nil {
		um.Id = bson.NewObjectId()
	}
	_, err = UserCollection(session).UpsertId(um.Id, um)
	if err != nil {
		return err
	}
	return nil
}
