package migration

import (
	r "github.com/dancannon/gorethink"
)

func CreateGames(db string, session *r.Session) error {
	_, err := r.Db(db).TableCreate("games").RunWrite(session)
	return err
}
