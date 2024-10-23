package config

import cmap "github.com/orcaman/concurrent-map/v2"

type IConfigMgr[T any] interface {
	Reload(tableName string, t T)
	Get(t T, id int) *T
	GetAll() []*T
}

type ConfigManager[T any] struct {
	configs cmap.ConcurrentMap[int, any]
}
