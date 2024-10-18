package persistence

import (
	"gameserver/common/logger"
	"reflect"
)

var repositories = map[any]interface{}{}

func RegisterRepository(model any, repo interface{}) {
	repositories[model] = repo
}

func Start() {
	logger.Logger.Info("Starting repositories...")
}

func Stop() {
	flushAllRepositories()
}

func flushAllRepositories() {
	logger.Logger.Info("Flushing all repositories...")
	for _, repo := range repositories {
		repoVal := reflect.ValueOf(repo)
		flushMethod := repoVal.MethodByName("Flush")
		if flushMethod.IsValid() && flushMethod.Type().NumIn() == 0 { // 确保Flush方法存在且无参数
			flushMethod.Call(nil) // 调用Flush方法
		}
	}
}
