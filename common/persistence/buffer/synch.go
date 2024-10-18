package buffer

import (
	"gorm.io/gorm"
)

type SyncBuffer[K string | int64, T any] struct {
	db *gorm.DB
}

func NewSyncBuffer[K string | int64, T any](db *gorm.DB) *SyncBuffer[K, T] {
	return &SyncBuffer[K, T]{
		db: db,
	}
}

// Add 方法实现
func (d *SyncBuffer[K, T]) Add(entity *T) *T {
	d.db.Create(entity)
	return entity
}

// Update 方法实现
func (d *SyncBuffer[K, T]) Update(entity *T) {
	d.db.Updates(entity)
}

// Remove 方法实现
func (d *SyncBuffer[K, T]) Remove(id K) {
	var entity T
	d.db.Model(entity).Where("id = ?", id).Delete(nil)
}

// RemoveAll 方法实现
func (d *SyncBuffer[K, T]) RemoveAll() {
	// 清空缓存并触发刷新
	var entity = new(T)
	d.db.Delete(&entity)
}

// Flush 方法实现
func (d *SyncBuffer[K, T]) Flush() {
}
