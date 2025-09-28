package config

import (
	"os"
)

type Config struct {
	Port       string
	SessionKey []byte
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	key := os.Getenv("SESSION_KEY")
	if key == "" {
		key = "dev-only-change-me"
	}

	return Config{Port: port, SessionKey: []byte(key)}
}

func (c Config) Addr() string {
	return ":" + c.Port
}
