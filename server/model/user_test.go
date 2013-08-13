package model

import (
	"testing"
)

func TestUserSavingAndLoading(t *testing.T) {
	if modelTestShouldRun() {
		cleanTestingDatabase()
		session, err := Connect()
		if err != nil {
			t.Fatal(err)
		}
		defer session.Close()
		um := UserModel{}
		um.Email = "fart@gmail.com"
		err = um.Save()
		if err != nil {
			t.Fatal(err)
		}
		if um.Id == nil {
			t.Fatal("User doesn't have an ID")
		}

		newUm, err := LoadUserByEmail("fart@gmail.com")
		if err != nil {
			t.Fatal(err)
		}
		if newUm.Id != um.Id {
			t.Fatal("ids don't match between old and new")
		}
		if newUm.Email != um.Email {
			t.Fatal("emails don't match between old and new")
		}
	}
}
