package log

import (
	"gameserver/common/database"
	"gameserver/common/persistence"
	"gameserver/common/persistence/repository"
	"gameserver/common/utils"
)

type DotLogin struct {
	ID         int64  `gorm:"column:id;primaryKey" `
	DayIndex   int    `gorm:"column:day_index;primaryKey" monthSharding:"true" partition:"day_index"`
	FirstTime  *int64 `gorm:"column:first_time"`
	LastTime   *int64 `gorm:"column:last_time" onupdate:"repeat"`
	TotalCount *int   `gorm:"column:total_count" onupdate:"total"`
}

func (d *DotLogin) TableName() string {
	//DayIndex格式为yyyyMMdd
	return utils.GetMonthTbName("dot_login", d.DayIndex)
}

//// BeforeCreate 回调函数
//func (d *DotLogin) BeforeCreate(tx *gorm.DB) (err error) {
//	// 判断表是否存在，如果不存在则创建表
//	tableName := d.TableName()
//	if !tx.Migrator().HasTable(tableName) {
//		err = tx.Migrator().CreateTable(&DotLogin{})
//		if err != nil {
//			return err
//		}
//	}
//	fmt.Println(fmt.Sprintf("BeforeCreate %s", tableName))
//	return
//}

var DotLoginRepository *repository.LoggerRepository[int64, DotLogin]

func init() {
	dotLogin := &DotLogin{}
	DotLoginRepository = repository.NewLoggerRepository[int64, DotLogin](database.GetLogDB(), "dot_login", true)
	persistence.RegisterRepository(dotLogin, DotLoginRepository)
}
