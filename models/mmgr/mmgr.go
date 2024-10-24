package mmgr

import (
	"container/list"
	"fmt"
	"gameserver/core/logger"
	"reflect"
)

var (
	// 注册所有数据库实体
	models       = list.New()
	repositories = list.New()
)

// 初始化所有实体的接口
type IModelInit interface {
	//反射调用InitRepository方法
	InitRepository()
}

// 注册实体和仓库
func RegisterModel(model any) {
	models.PushBack(model)
}
func RegisterRepository(repo any) {
	repositories.PushBack(repo)
	logger.Logger.Info(fmt.Sprintf("Registered repository for %v", repo))
}

func Start() {
	logger.Logger.Info("Starting repositories...")
	for i := models.Front(); i != nil; i = i.Next() {
		model := i.Value
		repoVal := reflect.ValueOf(model)
		initMethod := repoVal.MethodByName("InitRepository")
		if initMethod.IsValid() && initMethod.Type().NumIn() == 0 { // 确保InitRepository方法存在且无}
			initMethod.Call(nil)
			logger.Logger.Info(fmt.Sprintf("Initialized repository for %v", repoVal.Type()))
		}
	}
}

func Stop() {
	flushAllRepositories()
}

func flushAllRepositories() {
	logger.Logger.Info("Flushing all repositories...")
	for i := repositories.Front(); i != nil; i = i.Next() {
		model := i.Value
		repoVal := reflect.ValueOf(model)
		flushMethod := repoVal.MethodByName("Flush")
		if flushMethod.IsValid() && flushMethod.Type().NumIn() == 0 { // 确保Flush方法存在且无参数
			flushMethod.Call(nil) // 调用Flush方法
			logger.Logger.Info(fmt.Sprintf("Flushing repository for %v", repoVal.Type()))
		}
	}
}
