package redisclient

import (
	"context"
	"valeth-twice-management-api/internal/config"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func Init() {
    addr := config.Get("REDIS_ADDR", "localhost:6379")
    Rdb = redis.NewClient(&redis.Options{
        Addr: addr,
    })
}
