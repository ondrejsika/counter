package server

import (
	"github.com/ondrejsika/counter/internal/server"
)

type ServerConfig struct {
	DontRunMigrations bool
}

func Server(config ServerConfig) {
	server.Server(
		config.DontRunMigrations,
	)
}
