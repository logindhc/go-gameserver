package redis

import (
	"context"
	"gameserver/common/utils"
	"gameserver/conf"
	"gameserver/core/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

var (
	rdb            *redis.Client
	userIdRedisKey = "openId:"
	expTime        = time.Hour * 48
)

func InitRedis() {
	r := conf.GameConfig.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:         r.Addr,
		Password:     r.Password, // no password set
		DB:           r.DB,       // use default DB
		PoolSize:     r.PoolSize,
		MinIdleConns: r.MinIdleConns,
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	rdb = redisClient
	logger.Logger.Info("redis init success", zap.String("redis", redisClient.String()))
}

func GetRedisClient() *redis.Client {
	return rdb
}

func GetLoginToken(openId string) int64 {
	key := userIdRedisKey + openId
	userStr, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return utils.StrToInt64(userStr)
}

func SetLoginToken(openId string, userId int64) {
	key := userIdRedisKey + openId
	rdb.Set(context.Background(), key, userId, expTime)
}
