package main

import (
	"github.com/go-redis/redis"
	"github.com/haze518/zettel-bot/app"
)

func main() {
	client := new(zettel_bot.DBClient)
	*client = &zettel_bot.DropboxCLient{}
	storage := &zettel_bot.Storage{LifetimeSecond: 60}
	cl := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})
	storage.RedisClient = cl
	app := &zettel_bot.App{Client: client, Storage: storage}
	zettel_bot.Serve("5784759575:AAENlZ0UCZwhiifz8NgZwt6WqmFp2OmLCd8", app)
}
