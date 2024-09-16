package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
	Rdb *redis.Client
)

func InitRedis() *redis.Client {
	isTesting := os.Getenv("Testing")
	if isTesting == "false" {
		opts, err := redis.ParseURL(os.Getenv("REDIS_LOCALHOST"))
		if err != nil {
			panic(err)
		}
		Rdb = redis.NewClient(opts)
	} else {
		Rdb = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		})
	}
	err := Rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatal("err connecting tor redis, make sure docker img is running: ", err)
		return nil
	}
	return Rdb
}
