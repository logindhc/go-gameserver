package main

import (
	"gameserver/appserver/log"
	"gameserver/appserver/models"
	"gameserver/common/database"
)

var defaultModels = []interface{}{
	&models.Account{},
	&models.User{},
	&log.DotDevice{},
}

var logModels = []interface{}{
	&log.DotLogin{},
}

func init() {
	database.AutoMigrateByDbType(database.DEFAULT_DB, defaultModels)
	database.AutoMigrateByDbType(database.LOG_DB, logModels)
}
