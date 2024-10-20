package redis

import (
	"context"
	"fmt"
	"gameserver/common/logger"
	"gameserver/common/utils"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

var (
	name, cfgType, path = "redis", "yml", "."
	RedisConfig         *Config
	RDB                 *redis.Client
	userIdRedisKey      = "openId:"
	expTime             = time.Hour * 48
)

type Config struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

func initConfig() {
	viper.SetConfigName(name)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error(fmt.Sprintf("read %s config error: %v", name, err))
		return
	}
	redisConfig := &Config{}
	if err := viper.Unmarshal(redisConfig); err != nil {
		logger.Logger.Error(fmt.Sprintf("%s config unbale to decode into struct: %v", name, err))
		return
	}
	RedisConfig = redisConfig
	logger.Logger.Info(fmt.Sprintf("%s config init success", name))
}

func init() {
	initConfig()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     RedisConfig.Addr,
		Password: RedisConfig.Password, // no password set
		DB:       RedisConfig.DB,       // use default DB
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	RDB = redisClient
	logger.Logger.Info("redis init success", zap.String("redis", redisClient.String()))
}

func GetLoginToken(openId string) int64 {
	key := userIdRedisKey + openId
	userStr, err := RDB.Get(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return utils.StrToInt64(userStr)
}

func SetLoginToken(openId string, userId int64) {
	key := userIdRedisKey + openId
	RDB.Set(context.Background(), key, userId, expTime)
}
