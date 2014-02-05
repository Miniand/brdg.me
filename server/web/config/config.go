package config

import (
	"os"
)

const (
	SERVER_ADDRESS = "BRDGME_WEB_SERVER_ADDRESS"
)

var defaults = map[string]string{
	SERVER_ADDRESS: "127.0.0.1:9998",
}

func Get(key string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaults[key]
	}
	return v
}
