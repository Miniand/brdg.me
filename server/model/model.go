package model

import (
	"os"

	migrationDefinition "github.com/Miniand/brdg.me/server/model/migration"
	r "github.com/dancannon/gorethink"
)

var initialised = false

func DatabaseName() string {
	db := os.Getenv("BRDGME_DB_DATABASE")
	if db != "" {
		return db
	}
	return "brdgme"
}

func DatabaseAddr() string {
	addr := os.Getenv("BRDGME_DB_ADDRESS")
	if addr == "" {
		addr = "localhost:28015"
	}
	return addr
}

func Connect() (*r.Session, error) {
	if !initialised {
		initialised = true
		if err := initDb(); err != nil {
			initialised = false
			return nil, err
		}
	}
	return r.Connect(r.ConnectOpts{
		Address:  DatabaseAddr(),
		Database: DatabaseName(),
	})
}

type migration struct {
	Version string
	Up      func(db string, session *r.Session) error
}

var migrations = []migration{
	migration{"0001", migrationDefinition.CreateGames},
	migration{"0002", migrationDefinition.CreateUsers},
}

func migrate() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	for _, m := range migrations {
		hasRun, err := migrationHasRun(m.Version)
		if err != nil {
			return err
		}
		if !hasRun {
			if err := m.Up(DatabaseName(), session); err != nil {
				return err
			}
			if err := setMigrationHasRun(m.Version); err != nil {
				return err
			}
		}
	}
	return nil
}

func migrationHasRun(version string) (bool, error) {
	session, err := Connect()
	if err != nil {
		return false, err
	}
	res, err := r.Table("migrations").Filter(map[string]interface{}{
		"version": version,
	}).Run(session)
	if err != nil {
		return false, err
	}
	return !res.IsNil(), nil
}

func setMigrationHasRun(version string) error {
	session, err := Connect()
	if err != nil {
		return err
	}
	_, err = r.Table("migrations").Insert(map[string]interface{}{
		"version": version,
	}).RunWrite(session)
	return err
}

func createDb() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	r.DbCreate(DatabaseName()).RunWrite(session)
	r.Db(DatabaseName()).TableCreate("migrations").RunWrite(session)
	return nil
}

func dropDb() error {
	session, err := Connect()
	if err != nil {
		return err
	}
	r.DbDrop(DatabaseName()).RunWrite(session)
	return nil
}

func initDb() error {
	if err := createDb(); err != nil {
		return err
	}
	if err := migrate(); err != nil {
		return err
	}
	return nil
}
