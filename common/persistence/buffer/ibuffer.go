package buffer

import (
	"reflect"
)

type Buffer[K string | int64, T any] interface {
	Add(entity *T) *T
	Update(entity *T)
	Remove(id K)
	RemoveAll()
	Flush()
}

var flushIntervals = 3 // 默认的刷新间隔 3+(rand(1~3))分钟

// getKey 是一个辅助函数，用于从实体中提取键
func getKey[T any](entity *T) any {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		panic("ID field not found")
	}
	id, ok := idField.Interface().(any)
	if !ok {
		panic("ID is not an integer")
	}
	return id
}
