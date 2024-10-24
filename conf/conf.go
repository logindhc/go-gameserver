package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var GameConfig *config

type config struct {
	App       *app        `yaml:"app"`
	Log       *log        `yaml:"log"`
	Redis     *redis      `yaml:"redis"`
	Databases []*database `yaml:"databases"`
}

type app struct {
	Name     string `yaml:"name"`
	TcpPort  string `yaml:"tcpPort"`
	HttpPort string `yaml:"httpPort"`
	MisKey   string `yaml:"misKey"`
	OpenWS   bool   `yaml:"openWs"`
	JsonPath string `yaml:"jsonPath"`
}

type database struct {
	DbType       string        `yaml:"dbType"`
	DSN          string        `yaml:"dsn"`
	MaxIdleCount int           `yaml:"maxIdleCount"`
	MaxOpenCount int           `yaml:"maxOpenCount"`
	MaxLifetime  time.Duration `yaml:"maxLifetime"`
	AutoMigrate  bool          `yaml:"autoMigrate"`
}

type log struct {
	Level       string `yaml:"level"`
	FilePath    string `yaml:"filePath"`
	ErrFilePath string `yaml:"errFilePath"`
	MaxAge      int    `yaml:"maxAge"`
	MaxSize     int    `yaml:"maxSize"`
	MaxBackups  int    `yaml:"maxBackups"`
	TimeFormat  string `yaml:"timeFormat"`
}
type redis struct {
	Addr         string `yaml:"addr"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
}

func InitConfig(path string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path) // 添加搜索路径
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("InitConfig conf error: %v", err))
		return
	}
	temp := &config{}
	if err := viper.Unmarshal(temp); err != nil {
		panic(fmt.Sprintf("InitConfig unbale to decode into struct: %v", err))
		return
	}
	GameConfig = temp

	fmt.Println("InitConfig init success")
}
