package server

import (
	"github.com/ondrejsika/counter/internal/server"
)

type ServerOptions struct {
	DontRunMigrations bool
}

func Server(opts ServerOptions) {
	server.Server(
		opts.DontRunMigrations,
	)
}
