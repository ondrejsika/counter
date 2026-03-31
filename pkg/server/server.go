package server

import (
	"github.com/ondrejsika/counter/internal/server"
)

type ServerOptions struct {
	DontRunMigrations bool
	VersionOverride   string
}

func Server(opts ServerOptions) {
	server.Server(
		opts.DontRunMigrations,
		opts.VersionOverride,
	)
}
