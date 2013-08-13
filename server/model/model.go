package model

import (
	"labix.org/v2/mgo"
	"os"
)

func DatabaseName() string {
	db := os.Getenv("BRDGME_MONGODB_DATABASE")
	if db != "" {
		return db
	}
	return "brdgme"
}

func Connect() (*mgo.Session, error) {
	addr := os.Getenv("BRDGME_MONGODB_ADDRESS")
	if addr == "" {
		addr = "localhost"
	}
	return mgo.Dial(addr)
}
