package database

import (
	"fmt"
	"gameserver/common/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DEFAULT_DB, LOG_DB = "default", "log"

	name, cfgType, path, key = "database", "yml", "./conf", "databases"
	databases                = map[string]*gorm.DB{}
	Configs                  = map[string]*Config{}
)

type Config struct {
	DbType       string        `mapstructure:"db_type" json:"db_type" yaml:"db_type"`
	DSN          string        `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	MaxIdleCount int           `mapstructure:"max_idle_count" json:"max_idle_count" yaml:"max_idle_count"`
	MaxOpenCount int           `mapstructure:"max_open_count" json:"max_open_count" yaml:"max_open_count"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime" json:"max_lifetime" yaml:"max_lifetime"`
	AutoMigrate  bool          `mapstructure:"auto_migrate" json:"auto_migrate" yaml:"auto_migrate"`
}

func initConfig() {
	viper.SetConfigName(name)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error(fmt.Sprintf("read %s config error: %v", name, err))
		return
	}
	var configs []Config
	if err := viper.UnmarshalKey(key, &configs); err != nil {
		logger.Logger.Error(fmt.Sprintf("%s config unbale to decode into struct: %v", name, err))
		return
	}
	for _, config := range configs {
		Configs[config.DbType] = &config
	}
	logger.Logger.Info(fmt.Sprintf("%s config init success", name))
}

func init() {
	initConfig()
	for _, c := range Configs {
		databases[c.DbType] = connection(c)
	}
	logger.Logger.Info(fmt.Sprintf("databases connection len %d", len(databases)))
}

func GetDB(dbType string) *gorm.DB {
	if dbType == "" {
		dbType = DEFAULT_DB
	}
	return databases[dbType]
}

func GetDefaultDB() *gorm.DB {
	return GetDB(DEFAULT_DB)
}

func GetLogDB() *gorm.DB {
	return GetDB(LOG_DB)
}

func connection(c *Config) *gorm.DB {
	// 使用Zap来创建GORM的日志实现
	gormLogger := &logger.GormLoggerAdapter{Logger: logger.Logger.Sugar()}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s connection error %v", c, err))
		panic("connection db timeout")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s get db persistence error %v", c, err))
		panic("get db persistence error")
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleCount)
	sqlDB.SetMaxOpenConns(c.MaxOpenCount)
	sqlDB.SetConnMaxLifetime(c.MaxLifetime)
	return db
}

func AutoMigrateByDbType(dbType string, models []interface{}) {
	autoMigrate := Configs[dbType].AutoMigrate
	if !autoMigrate {
		return
	}
	db := GetDB(dbType)
	for _, model := range models {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model)
		if err != nil {
			return
		}
	}
	logger.Logger.Info(fmt.Sprintf("persistence autoMigrate success auto -> %v dbType: %s log: %d", autoMigrate, dbType, len(models)))
}
