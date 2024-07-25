package main

import (
	"log"
	"os"

	"github.com/ondrejsika/counter/backend_inmemory"
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
	} else {
		log.Fatalf(`no backend "%s" exists, you can use "redis" (default) or "inmemory"\n`, backend)
	}

	server.Server(
		doCountFunc,
		getCountFunc,
	)
}
