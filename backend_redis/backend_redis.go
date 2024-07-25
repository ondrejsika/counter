package backend_redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func DoCountRedis(redisHost, hostname string) (int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

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

func GetCountRedis(redisHost, hostname string) (int, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379",
		Password: "",
		DB:       0,
	})

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
