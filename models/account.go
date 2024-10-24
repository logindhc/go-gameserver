package models

import (
	"gameserver/common/persistence/repository"
	"gameserver/core/database"
	"gameserver/models/mmgr"
)

type Account struct {
	ID            string `gorm:"primaryKey"`
	OpenId        string `gorm:"index"`
	Channel       int    `gorm:"index"`
	ChannelOpenId string
	Platform      int
	RegisterTime  int64
	TotalLoginDay int
}

func (a *Account) TableName() string {
	return "account"
}

var AccountRepository *repository.LruRepository[string, Account]

func init() {
	mmgr.RegisterModel(&Account{})
}
func (a *Account) InitRepository() {
	AccountRepository = repository.NewLruRepository[string, Account](database.GetDefaultDB(), a.TableName())
	mmgr.RegisterRepository(AccountRepository)
}
