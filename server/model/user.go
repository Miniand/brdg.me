package model

import (
	r "github.com/dancannon/gorethink"
)

type UserModel struct {
	Id           string `gorethink:"id,omitempty"`
	Email        string
	Unsubscribed bool
}

func UserTable() r.RqlTerm {
	return r.Table("users")
}

func LoadUser(id string) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	row, err := UserTable().Get(id).RunRow(session)
	if err != nil {
		return nil, err
	}
	m := &UserModel{}
	err = row.Scan(m)
	return m, err
}

func FirstUserByEmail(email string) (*UserModel, error) {
	session, err := Connect()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	row, err := UserTable().Filter(map[string]interface{}{
		"Email": email,
	}).RunRow(session)
	if err != nil {
		return nil, err
	}
	if row.IsNil() {
		return nil, nil
	}
	m := &UserModel{}
	err = row.Scan(m)
	return m, err
}

func (um *UserModel) Save() error {
	var rqlTerm r.RqlTerm
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
