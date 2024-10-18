package cache

import (
	"github.com/redis/go-redis/v9"
)

type RedisCache[K string | int64, V any] struct {
	*redis.Client
}
