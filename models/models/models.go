package models

import (
	"gameserver/core/database"
	"gameserver/models"
)

var (
	// 配置需要自动更新数据库表配置的注册列表
	// 默认表
	defaultModels = []interface{}{
		&models.Account{},
		&models.User{},
	}
	// 日志表
	logModels = []interface{}{
		&models.DotLogin{},
		&models.DotDevice{},
	}
)

func InitAutoMigrate() {
	database.AutoMigrateByDbType(database.DEFAULT_DB, defaultModels)
	database.AutoMigrateByDbType(database.LOG_DB, logModels)
}
