package model

import (
	r "github.com/dancannon/gorethink"
)

type UserModel struct {
	Id           string `gorethink:"id,omitempty"`
	Email        string
	Unsubscribed bool
}

func UserTable() r.Term {
	return r.Table("users")
}

func LoadUser(id string) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	res, err := UserTable().Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	m := &UserModel{}
	err = res.One(m)
	return m, err
}

func FirstUserByEmail(email string) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	res, err := UserTable().Filter(map[string]interface{}{
		"Email": email,
	}).Run(session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return nil, nil
	}
	m := &UserModel{}
	err = res.One(m)
	return m, err
}

func (um *UserModel) Save() error {
	var rqlTerm r.Term
	session, err := Connect()
	if err != nil {
		return err
	}
	defer session.Close()
	if um.Id == "" {
		rqlTerm = UserTable().Insert(um)
	} else {
		rqlTerm = UserTable().Get(um.Id).Update(um)
	}
	res, err := rqlTerm.RunWrite(session)
	if err != nil {
		return err
	}
	if um.Id == "" {
		um.Id = res.GeneratedKeys[0]
	}
	return nil
}
