package buffer

import (
	"gorm.io/gorm"
)

type ASyncBuffer[K string | int64, T any] struct {
	Db *gorm.DB
}

func NewASyncBuffer[K string | int64, T any](db *gorm.DB) *ASyncBuffer[K, T] {
	return &ASyncBuffer[K, T]{
		Db: db,
	}
}

// Add 方法实现
func (d *ASyncBuffer[K, T]) Add(entity *T) *T {
	go func() {
		d.Db.Create(entity)
	}()
	return entity
}

// Update 方法实现
func (d *ASyncBuffer[K, T]) Update(entity *T) {
	go func() {
		d.Db.Updates(entity)
	}()
}

// Remove 方法实现
func (d *ASyncBuffer[K, T]) Remove(id K) {
	var entity T
	go func() {
		d.Db.Delete(&entity, id)
	}()
}

// RemoveAll 方法实现
func (d *ASyncBuffer[K, T]) RemoveAll() {
	// 清空缓存并触发刷新
	var entity = new(T)
	go func() {
		d.Db.Delete(&entity)
	}()
}

// Flush 方法实现
func (d *ASyncBuffer[K, T]) Flush() {
}
