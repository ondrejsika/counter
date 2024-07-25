package main

import (
	"os"

	"github.com/ondrejsika/counter/backend_redis"
	"github.com/ondrejsika/counter/server"
)

func main() {
	hostname, _ := os.Hostname()

	redisHost := "127.0.0.1"
	envRedisHost := os.Getenv("REDIS")
	if envRedisHost != "" {
		redisHost = envRedisHost
	}

	server.Server(
		// do count
		func() (int, error) {
			return backend_redis.DoCountRedis(redisHost, hostname)
		},
		// get count
		func() (int, error) {
			return backend_redis.GetCountRedis(redisHost, hostname)
		},
	)
}
