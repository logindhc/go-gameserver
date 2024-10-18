package models

import (
	"gameserver/common/database"
	"gameserver/common/persistence"
	"gameserver/common/persistence/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	account := &Account{}
	AccountRepository = repository.NewLruRepository[string, Account](database.GetDefaultDB(), account.TableName())
	persistence.RegisterRepository(account, AccountRepository)
}
func log() {
	database.GetDefaultDB().Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "ID"}},
		DoUpdates: clause.AssignmentColumns([]string{"ID"})}).Exec("")

	accounts := []Account{}
	database.GetDefaultDB().Session(&gorm.Session{FullSaveAssociations: true}).CreateInBatches(accounts, 100)
}
