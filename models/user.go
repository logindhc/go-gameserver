package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gameserver/common/persistence/repository"
	"gameserver/core/database"
	"gameserver/models/mmgr"
)

type User struct {
	ID            int64  `gorm:"primaryKey"`
	OpenId        string `gorm:"index"`
	Platform      int    `gorm:"index"`
	Channel       int    `gorm:"index"`
	RegisterTime  int64
	Device        string
	Country       int
	LastLoginIp   string
	LastLoginTime int64
	LastLoginDay  int
	TotalLoginDay int
	Level         int
	Items         ItemMap `gorm:"type:longtext"`
}

func (u *User) GetItems() ItemMap {
	if u.Items == nil {
		return ItemMap{}
	}
	return u.Items
}

type ItemMap map[int]int

func (im ItemMap) Value() (driver.Value, error) {
	if im == nil {
		return nil, nil
	}
	data, err := json.Marshal(im)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (im *ItemMap) Scan(value interface{}) error {
	if value == nil {
		*im = ItemMap{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var m ItemMap
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	*im = m
	return nil
}

func (u *User) TableName() string {
	return "user"
}

var UserRepository *repository.LruRepository[int64, User]

func init() {
	mmgr.RegisterModel(&User{})
}
func (u *User) InitRepository() {
	UserRepository = repository.NewLruRepository[int64, User](database.GetDefaultDB(), u.TableName())
	mmgr.RegisterRepository(UserRepository)
}
