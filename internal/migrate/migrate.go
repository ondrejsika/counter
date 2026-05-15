package migrate

import (
	"os"

	"github.com/ondrejsika/counter/internal/backend_postgres"
	"github.com/ondrejsika/counter/internal/backend_redis"
)

func Migrate() {
	var runMigrationsFunc func() error

	hostname, _ := os.Hostname()

	backend := "redis"
	envBackend := os.Getenv("BACKEND")
	if envBackend != "" {
		backend = envBackend
	}

	if backend == "redis" {
		redisHost := os.Getenv("REDIS")
		if redisHost == "" {
			redisHost = "127.0.0.1"
		}
		redisPassword := os.Getenv("REDIS_PASSWORD")
		redisSchema := os.Getenv("REDIS_SCHEMA")
		runMigrationsFunc = func() error {
			return backend_redis.RunFakeMigrations(redisHost, redisPassword, hostname, redisSchema)
		}
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
