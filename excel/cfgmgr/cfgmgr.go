package cfgmgr

import (
	"fmt"
	"gameserver/conf"
	"gameserver/core/logger"
	"gameserver/excel/excels"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"reflect"
	"sync/atomic"
	"time"
)

type ICfgManager interface {
	CfgName() string
	Load()
}

type CfgManager struct {
	update atomic.Bool
	timer  *time.Timer
}

func NewManager() *CfgManager {
	cm := &CfgManager{
		update: atomic.Bool{},
	}
	cm.initConfig()
	cm.load()
	return cm
}
func (c *CfgManager) initConfig() {
	viper.SetConfigFile(conf.GameConfig.App.JsonPath)
	viper.OnConfigChange(func(e fsnotify.Event) {
		//实际会触发两次
		c.update.Store(true)
	})
	viper.WatchConfig()
	c.startUpdateChecker()
	logger.Logger.Info(fmt.Sprintf("conf manager init success. len %v", len(excels.GetTables())))
}

// 启动定时检查更新
func (c *CfgManager) startUpdateChecker() {
	checkInterval := 5 * time.Second // 定时检查间隔时间
	c.timer = time.NewTimer(checkInterval)
	go func() {
		for {
			select {
			case <-c.timer.C:
				if c.update.Load() {
					c.load()
					c.update.Store(false)
				}
				c.timer.Reset(checkInterval)
			}
		}
	}()
}

func (c *CfgManager) load() {
	for name, cfg := range excels.GetTables() {
		cfgVal := reflect.ValueOf(cfg)
		flushMethod := cfgVal.MethodByName("Load")
		if flushMethod.IsValid() && flushMethod.Type().NumIn() == 0 {
			flushMethod.Call(nil)
			logger.Logger.Info(fmt.Sprintf("load conf %v succeed", name))
		}
	}
}
