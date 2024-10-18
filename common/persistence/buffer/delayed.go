package buffer

import (
	"fmt"
	"gameserver/common/concurrent"
	"gameserver/common/logger"
	"gameserver/common/persistence/cache"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type DelayedBuffer[K string | int64, T any] struct {
	cache      *cache.LRUCache[K, T]
	db         *gorm.DB
	prefix     string
	updates    *concurrent.ConcurrentSet[K]
	deletes    *concurrent.ConcurrentSet[K]
	bufferSize int
}

func NewDelayedBuffer[K string | int64, T any](db *gorm.DB, cache *cache.LRUCache[K, T], prefix string) *DelayedBuffer[K, T] {
	bufferSize := 100
	buffer := &DelayedBuffer[K, T]{
		cache:      cache,
		db:         db,
		prefix:     prefix,
		updates:    concurrent.NewConcurrentSet[K](),
		deletes:    concurrent.NewConcurrentSet[K](),
		bufferSize: bufferSize,
	}
	go buffer.flushLoop() // 启动后台任务处理更新与删除
	return buffer
}

// flushLoop 是一个后台循环，用于定期将缓存中的更改同步到数据库
func (d *DelayedBuffer[K, T]) flushLoop() {
	interval := time.Duration(flushIntervals+rand.Intn(flushIntervals)) * time.Minute
	logger.Logger.Info(fmt.Sprintf("%s# start flushLoop task interval %d", d.prefix, interval))
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.Flush()
		}

	}
}

// Add 方法实现
func (d *DelayedBuffer[K, T]) Add(entity *T) *T {
	k := getKey(entity)
	d.deletes.Remove(k.(K))
	tx := d.db.Create(entity)
	if tx.Error == nil {
		logger.Logger.Info(fmt.Sprintf("%s#id:%v 添加成功", d.prefix, k))
	} else {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v 添加失败", d.prefix, k))
	}
	return entity
}

// Update 方法实现
func (d *DelayedBuffer[K, T]) Update(entity *T) {
	id := getKey(entity)
	if d.deletes.Has(id.(K)) {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v 更新时已经被删除", d.prefix, id))
		return
	}
	d.updates.Add(id.(K))
	//fmt.Printf("updates %p %+d \n", d.updates, d.updates.Size())
	if d.updates.Size() >= d.bufferSize {
		d.Flush()
	}
}

// Remove 方法实现
func (d *DelayedBuffer[K, T]) Remove(id K) {
	d.deletes.Add(id)
	d.updates.Remove(id)
	var entity T
	d.db.Model(entity).Where("id = ?", id).Delete(nil)
	d.deletes.Remove(id)
	logger.Logger.Info(fmt.Sprintf("%s#id:%v 删除成功", d.prefix, id))
}

// RemoveAll 方法实现
func (d *DelayedBuffer[K, T]) RemoveAll() {
	// 清空缓存并触发刷新
	d.cache.Clear()
	d.deletes.Clear()
	d.updates.Clear()
	var entity = new(*T)
	d.db.Delete(entity)
}

// Flush 方法实现
func (d *DelayedBuffer[K, T]) Flush() {
	// 处理更新
	size := d.updates.Size()
	if size <= 0 {
		return
	}
	logger.Logger.Info(fmt.Sprintf("%s# update num %d", d.prefix, size))
	all := d.updates.All()
	for _, id := range all {
		entity := d.cache.Get(id)
		if entity == nil {
			logger.Logger.Error(fmt.Sprintf("%s#id:%v 更新失败，缓存中不存在", d.prefix, id))
			continue
		}
		tx := d.db.Updates(entity)
		if tx.Error != nil {
			logger.Logger.Error(fmt.Sprintf("%s#id:%v 更新失败", d.prefix, id))
			continue
		}
		d.updates.Remove(id)
	}
}
