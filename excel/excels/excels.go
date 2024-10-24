package excels

import "gameserver/excel"

var (
	// 配置管理器，注册配置表，便于加载或更新
	// 先手动配置，后面自动生成
	tables = map[string]interface{}{
		"LevelConfig": &excel.LevelCfg{},
	}
)

func GetTables() map[string]interface{} {
	return tables
}
