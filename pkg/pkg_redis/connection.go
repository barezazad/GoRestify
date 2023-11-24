package pkg_redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Ctx context for redis commands
var Ctx = context.Background()

// DisplayLog to print redis results
var DisplayLog bool

// ConnectRedis create a connection to local redis it uses to cache
func ConnectRedis(address string, displayLog bool) (redisCon RedisCon, err error) {

	DisplayLog = displayLog

	opt, err := redis.ParseURL(address)
	if err != nil {
		return
	}

	redisCon.client = redis.NewClient(opt)
	var ping string
	if ping, err = redisCon.client.Ping(Ctx).Result(); err != nil {
		return
	}

	if DisplayLog {
		fmt.Println("PONG:", ping)
	}

	return
}
