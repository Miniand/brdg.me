package migration

import (
	r "github.com/dancannon/gorethink"
)

func CreateGames(db string, session *r.Session) error {
	_, err := r.DB(db).TableCreate("games").RunWrite(session)
	return err
}
