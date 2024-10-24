package excel

import (
	"fmt"
	"gameserver/conf"
	"gameserver/core/logger"
	"github.com/spf13/viper"
	"path/filepath"
)

var LevelConfig = map[int]LevelCfg{}

type LevelCfg struct {
	Id         int    `json:"id"`
	Reward     []int  `json:"Reward"`
	LostReward string `json:"LostReward"`
}

func (l *LevelCfg) CfgName() string {
	return "LevelConfig"
}

func (l *LevelCfg) Load() {
	var temp map[int]LevelCfg
	err := loadConfig(l.CfgName(), &temp)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s load error %v", l.CfgName(), err))
		return
	}
	LevelConfig = temp
}

// 通用配置文件加载函数
func loadConfig(configName string, target interface{}) error {
	configName = fmt.Sprintf("%s.json", configName)
	viper.SetConfigFile(filepath.Join(conf.GameConfig.App.JsonPath, configName))
	if err := viper.ReadInConfig(); err != nil {
		logger.Logger.Error(fmt.Sprintf("read %s excel error: %v", configName, err))
		return err
	}
	if err := viper.Unmarshal(target); err != nil {
		logger.Logger.Error(fmt.Sprintf("unmarshal %s excel error: %v", configName, err))
		return err
	}
	return nil
}
