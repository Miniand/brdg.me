package migration

import (
	r "github.com/dancannon/gorethink"
)

const (
	AuthTokenTable = "authtokens"
)

func CreateAuthtokens(db string, session *r.Session) error {
	if _, err := r.Db(db).TableCreate(AuthTokenTable).
		RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table(AuthTokenTable).IndexCreate("Token").
		RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table(AuthTokenTable).IndexWait("Token").
		Run(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table(AuthTokenTable).IndexCreateFunc("UserId:Token",
		func(row r.Term) interface{} {
			return []interface{}{
				row.Field("UserId"),
				row.Field("Token"),
			}
		}).RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table(AuthTokenTable).IndexWait("UserId:Token").
		Run(session); err != nil {
		return err
	}
	return nil
}
