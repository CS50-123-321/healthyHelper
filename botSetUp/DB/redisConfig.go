package DB

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
		opts, err := redis.ParseURL(os.Getenv("REDIS_FLY"))
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
		log.Fatal("err: ", err)
		return nil
	}
	return Rdb
}
