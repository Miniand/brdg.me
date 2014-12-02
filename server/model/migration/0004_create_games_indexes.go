package migration

import r "github.com/dancannon/gorethink"

func CreateGamesIndexes(db string, session *r.Session) error {
	// Index by IsFinished and WhoseTurn for current turn search
	if _, err := r.Db(db).Table("games").IndexCreateFunc(
		"IsFinished:WhoseTurn",
		func(row r.Term) interface{} {
			return row.Field("WhoseTurn").Map(func(wt r.Term) interface{} {
				return []interface{}{row.Field("IsFinished"), wt}
			})
		},
		r.IndexCreateOpts{Multi: true},
	).RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table("games").IndexWait("IsFinished:WhoseTurn").
		Run(session); err != nil {
		return err
	}
	// Index by FinishedAt for sorting
	if _, err := r.Db(db).Table("games").IndexCreate("FinishedAt").
		Run(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table("games").IndexWait("FinishedAt").
		Run(session); err != nil {
		return err
	}
	// Index by IsFinished and PlayerList to get recently finished and active
	if _, err := r.Db(db).Table("games").IndexCreateFunc(
		"IsFinished:PlayerList",
		func(row r.Term) interface{} {
			return row.Field("PlayerList").Map(func(pl r.Term) interface{} {
				return []interface{}{row.Field("IsFinished"), pl}
			})
		},
		r.IndexCreateOpts{Multi: true},
	).RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table("games").IndexWait("IsFinished:PlayerList").
		Run(session); err != nil {
		return err
	}
	// Index by PlayerList to get all games for a player
	if _, err := r.Db(db).Table("games").IndexCreate(
		"PlayerList",
		r.IndexCreateOpts{Multi: true},
	).RunWrite(session); err != nil {
		return err
	}
	if _, err := r.Db(db).Table("games").IndexWait("PlayerList").
		Run(session); err != nil {
		return err
	}
	return nil
}
