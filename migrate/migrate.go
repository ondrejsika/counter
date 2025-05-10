package migrate

import (
	_ "embed"
	"os"

	"github.com/ondrejsika/counter/backend_postgres"
)

func Migrate() {
	var runMigrationsFunc func() error

	hostname, _ := os.Hostname()

	backend := "redis"
	envBackend := os.Getenv("BACKEND")
	if envBackend != "" {
		backend = envBackend
	}

	if backend == "postgres" {
		postgresHost := "127.0.0.1"
		envPostgresHost := os.Getenv("POSTGRES_HOST")
		if envPostgresHost != "" {
			postgresHost = envPostgresHost
		}

		postgresUser := "postgres"
		envPostgresUser := os.Getenv("POSTGRES_USER")
		if envPostgresUser != "" {
			postgresUser = envPostgresUser
		}

		postgresPassword := "pg"
		envPostgresPassword := os.Getenv("POSTGRES_PASSWORD")
		if envPostgresPassword != "" {
			postgresPassword = envPostgresPassword
		}

		postgresDatabase := "postgres"
		envPostgresDatabase := os.Getenv("POSTGRES_DATABASE")
		if envPostgresDatabase != "" {
			postgresDatabase = envPostgresDatabase
		}

		postgresSslmode := "disable"
		envPostgresSslmode := os.Getenv("POSTGRES_SSLMODE")
		if envPostgresSslmode != "" {
			postgresSslmode = envPostgresSslmode
		}

		runMigrationsFunc = func() error {
			return backend_postgres.RunMigrations(
				postgresHost, 5432, postgresUser, postgresPassword, postgresDatabase, postgresSslmode, hostname,
			)
		}
	}

		runMigrationsFunc()
}
