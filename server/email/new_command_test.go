package main

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	nc := NewCommand{}
	result := nc.Parse(`	   new  texas_holdem beefsack@gmail.com    
    kshaushau@gmail.com 	baconheist@gmail.com striker203@gmail.com`)
	if result == nil {
		t.Fatal("Couldn't find command")
	}
	if len(result) < 2 {
		t.Fatal("Couldn't find game name and emails")
	}
	if result[1] != "texas_holdem" {
		t.Fatal("Argument 1 isn't texas_holdem")
	}
	if result[2] != ` beefsack@gmail.com    
    kshaushau@gmail.com 	baconheist@gmail.com striker203@gmail.com` {
		t.Fatal("Argument 2 isn't the list of email addresses")
	}
}
