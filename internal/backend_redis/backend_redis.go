package backend_redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func DoCountRedis(redisHost, redisPassword, hostname string) (int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: redisPassword,
		DB:       0,
	})
	defer rdb.Close()

	if os.Getenv("REDIS_SCHEMA") != "" {
		schemaVal, err := rdb.Get(ctx, "schema").Result()
		if err != nil {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
			return -1, err
		}
		dbSchema, err := strconv.Atoi(schemaVal)
		if err != nil {
			return -1, fmt.Errorf("invalid db schema version: %s", schemaVal)
		}
		envSchema, err := strconv.Atoi(os.Getenv("REDIS_SCHEMA"))
		if err != nil {
			return -1, fmt.Errorf("invalid REDIS_SCHEMA env value: %s", os.Getenv("REDIS_SCHEMA"))
		}
		if dbSchema < envSchema {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("schema too old: db=%d required=%d", dbSchema, envSchema))
			return -1, fmt.Errorf("schema too old: db=%d required=%d", dbSchema, envSchema)
		}
	}

	val, err := rdb.Get(ctx, "counter").Result()
	if err != nil {
		if err == redis.Nil {
			val = "0"
		} else {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
			return -1, err
		}
	}

	counter, err := strconv.Atoi(val)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	err = rdb.Set(ctx, "counter", strconv.Itoa(counter+1), 0).Err()
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	return counter + 1, nil
}

func GetCountRedis(redisHost, redisPassword, hostname string) (int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: redisPassword,
		DB:       0,
	})
	defer rdb.Close()

	if os.Getenv("REDIS_SCHEMA") != "" {
		schemaVal, err := rdb.Get(ctx, "schema").Result()
		if err != nil {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
			return -1, err
		}
		dbSchema, err := strconv.Atoi(schemaVal)
		if err != nil {
			return -1, fmt.Errorf("invalid db schema version: %s", schemaVal)
		}
		envSchema, err := strconv.Atoi(os.Getenv("REDIS_SCHEMA"))
		if err != nil {
			return -1, fmt.Errorf("invalid REDIS_SCHEMA env value: %s", os.Getenv("REDIS_SCHEMA"))
		}
		if dbSchema < envSchema {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("schema too old: db=%d required=%d", dbSchema, envSchema))
			return -1, fmt.Errorf("schema too old: db=%d required=%d", dbSchema, envSchema)
		}
	}

	val, err := rdb.Get(ctx, "counter").Result()
	if err != nil {
		if err == redis.Nil {
			val = "0"
		} else {
			log.Error().
				Str("hostname", hostname).
				Msg(fmt.Sprintf("error=%s", err))
			return -1, err
		}
	}

	counter, err := strconv.Atoi(val)
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	return counter, nil
}

func RunFakeMigrations(redisHost, redisPassword, hostname, schema string) error {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: redisPassword,
		DB:       0,
	})
	defer rdb.Close()

	err := rdb.Set(ctx, "schema", schema, 0).Err()
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return err
	}

	return nil
}
