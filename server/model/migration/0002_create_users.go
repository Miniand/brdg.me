package migration

import (
	r "github.com/dancannon/gorethink"
)

func CreateUsers(db string, session *r.Session) error {
	_, err := r.DB(db).TableCreate("users").RunWrite(session)
	return err
}
