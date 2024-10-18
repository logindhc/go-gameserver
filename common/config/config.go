package config

import (
	"fmt"
	"gameserver/common/logger"
	"github.com/spf13/viper"
)

var (
	name, cfgType, path = "server", "yml", "."
	ServerConfig        *Config
)

type Config struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	TcpPort  string `mapstructure:"tcp_port" json:"tcp_port" yaml:"tcp_port"`
	HttpPort string `mapstructure:"http_port" json:"http_port" yaml:"http_port"`
	MisKey   string `mapstructure:"mis_key" json:"mis_key" yaml:"mis_key"`
}

func init() {
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
	ServerConfig = redisConfig
	logger.Logger.Info(fmt.Sprintf("%s config init success", name))
}
