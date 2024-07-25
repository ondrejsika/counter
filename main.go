package main

import (
	"log"
	"os"

	"github.com/ondrejsika/counter/backend_inmemory"
	"github.com/ondrejsika/counter/backend_postgres"
	"github.com/ondrejsika/counter/backend_redis"
	"github.com/ondrejsika/counter/server"
)

func main() {
	hostname, _ := os.Hostname()

	backend := "redis"
	envBackend := os.Getenv("BACKEND")
	if envBackend != "" {
		backend = envBackend
	}

	var doCountFunc func() (int, error)
	var getCountFunc func() (int, error)

	if backend == "redis" {
		redisHost := "127.0.0.1"
		envRedisHost := os.Getenv("REDIS")
		if envRedisHost != "" {
			redisHost = envRedisHost
		}
		doCountFunc = func() (int, error) { return backend_redis.DoCountRedis(redisHost, hostname) }
		getCountFunc = func() (int, error) { return backend_redis.GetCountRedis(redisHost, hostname) }
	} else if backend == "inmemory" {
		doCountFunc = func() (int, error) { return backend_inmemory.DoCountInMemory() }
		getCountFunc = func() (int, error) { return backend_inmemory.GetCountInMemory() }
	} else if backend == "postgres" {
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

		doCountFunc = func() (int, error) {
			return backend_postgres.DoCountPostgres(
				postgresHost, 5432, postgresUser, postgresPassword, postgresDatabase, hostname,
			)
		}
		getCountFunc = func() (int, error) {
			return backend_postgres.GetCountPostgres(
				postgresHost, 5432, postgresUser, postgresPassword, postgresDatabase, hostname,
			)
		}
	} else {
		log.Fatalf(`no backend "%s" exists, you can use "redis" (default), "postgres", or "inmemory"\n`, backend)
	}

	server.Server(
		doCountFunc,
		getCountFunc,
	)
}
