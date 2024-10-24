package database

import (
	"fmt"
	"gameserver/conf"
	"gameserver/core/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DEFAULT_DB, LOG_DB = "default", "log"
	databases          = map[string]*gorm.DB{}
)

func InitDatabase() {
	for _, c := range conf.GameConfig.Databases {
		databases[c.DbType] = connection(c.DbType, c.DSN, c.MaxLifetime, c.MaxOpenCount, c.MaxIdleCount)
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

func connection(dbtype string, dsn string, maxLifetime time.Duration, maxOpenCount int, maxIdleCount int) *gorm.DB {
	//使用Zap来创建GORM的日志实现
	gormLogger := &logger.GormLoggerAdapter{Logger: logger.Logger.Sugar()}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%v connection error %v", dbtype, err))
		panic("connection db timeout")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%v get db persistence error %v", dbtype, err))
		panic("get db persistence error")
	}
	sqlDB.SetMaxIdleConns(maxIdleCount)
	sqlDB.SetMaxOpenConns(maxOpenCount)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	return db
}

func AutoMigrateByDbType(dbType string, models []interface{}) {
	cfg := conf.GameConfig.Databases
	cfgType := false
	for _, c := range cfg {
		if c.DbType == dbType {
			cfgType = true
			break
		}
	}
	if !cfgType {
		return
	}
	db := GetDB(dbType)
	for _, model := range models {
		err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(model)
		if err != nil {
			return
		}
	}
	logger.Logger.Info(fmt.Sprintf("persistence autoMigrate success auto -> %v dbType: %s log: %d", cfgType, dbType, len(models)))
}
