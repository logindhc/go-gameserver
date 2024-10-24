package models

import (
	"gameserver/common/persistence/repository"
	"gameserver/core/database"
	"gameserver/models/mmgr"
)

type DotDevice struct {
	ID       string `gorm:"column:id;primaryKey" `
	Platform int    `gorm:"column:platform;primaryKey"`
	Channel  *int64 `gorm:"column:channel"`
	FChannel *int64 `gorm:"column:f_channel"`
	Country  *int64 `gorm:"column:country"`
	Time     *int64 `gorm:"column:time"`
}

func (log *DotDevice) TableName() string {
	return "dot_device"
}

var DotDeviceRepository *repository.LoggerRepository[string, DotDevice]

func init() {
	mmgr.RegisterModel(&DotDevice{})
}
func (log *DotDevice) InitRepository() {
	DotDeviceRepository = repository.NewLoggerRepository[string, DotDevice](database.GetLogDB(), "dot_device", false)
	mmgr.RegisterRepository(DotDeviceRepository)
}
