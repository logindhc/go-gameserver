package log

import (
	"gameserver/common/database"
	"gameserver/common/persistence"
	"gameserver/common/persistence/repository"
)

type DotDevice struct {
	ID       string `gorm:"column:id;primaryKey" `
	Platform int    `gorm:"column:platform;primaryKey"`
	Channel  *int64 `gorm:"column:channel"`
	FChannel *int64 `gorm:"column:f_channel"`
	Country  *int64 `gorm:"column:country"`
	Time     *int64 `gorm:"column:time"`
}

func (d *DotDevice) TableName() string {
	return "dot_device"
}

var DotDeviceRepository *repository.SynchRepository[string, DotDevice]

func init() {
	dotDevice := &DotDevice{}
	DotDeviceRepository = repository.NewSynchRepository[string, DotDevice](database.GetDefaultDB())
	persistence.RegisterRepository(dotDevice, DotDeviceRepository)
}
