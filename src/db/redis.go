package db

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})
