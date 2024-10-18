package repository

import (
	"fmt"
	"gameserver/common/logger"
	"gameserver/common/persistence/buffer"
	"gameserver/common/persistence/cache"
	"gorm.io/gorm"
	"reflect"
	"sync"
	"time"
)

type LruRepository[K string | int64, T any] struct {
	cache  *cache.LRUCache[K, T]
	db     *gorm.DB
	buffer *buffer.DelayedBuffer[K, T]
	prefix string
	lock   *sync.Mutex
}

func NewLruRepository[K string | int64, T any](db *gorm.DB, prefix string) *LruRepository[K, T] {
	// lru 过期设置为2个小时
	lruCache := cache.NewLRUCache[K, T](2*time.Hour, 2*time.Hour)
	r := &LruRepository[K, T]{
		db:     db,
		cache:  lruCache,
		buffer: buffer.NewDelayedBuffer[K, T](db, lruCache, prefix),
		prefix: prefix,
		lock:   &sync.Mutex{},
	}
	return r
}

func (r *LruRepository[K, T]) Get(id K) *T {
	// 从本地缓存中获取
	entity := r.cache.Get(id)
	if entity != nil {
		return entity
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	// 如果缓存中没有，则从数据库中获取
	tx := r.db.Where("id = ?", id).Find(&entity)

	if tx.RowsAffected == 0 {
		return nil
	}
	if entity != nil {
		e := r.cache.PutIfAbsent(id, entity)
		if e != nil {
			return e
		}
	}
	return entity
}

func (r *LruRepository[K, T]) GetAll() []*T {
	var entities []*T
	tx := r.db.Find(&entities)
	if tx.Error != nil {
		logger.Logger.Error(fmt.Sprintf("%s#all查询失败", r.prefix))
		return nil
	}
	return entities
}

func (r *LruRepository[K, T]) GetOrCreate(id K) *T {
	entity := r.Get(id)
	if entity == nil {
		entity = r.cache.Get(id)
		if entity == nil {
			entity = new(T)
			r.setId(entity, id)
			entity = r.Add(entity)
		}
	}
	return entity
}

func (r *LruRepository[K, T]) Add(entity *T) *T {
	if entity == nil {
		return nil
	}
	id := r.getId(entity)
	prev := r.cache.Get(id)
	if prev != nil {
		return prev
	}

	return r.buffer.Add(entity)
}

func (r *LruRepository[K, T]) Remove(id K) {
	r.cache.Remove(id)
	r.buffer.Remove(id)
}

func (r *LruRepository[K, T]) Update(entity *T) {
	id := r.getId(entity)
	r.cache.Put(id, entity)
	r.buffer.Update(entity)
}
func (r *LruRepository[K, T]) Flush() {
	r.buffer.Flush()
}

func (r *LruRepository[K, T]) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return r.db.Where(query, args...)
}

func (r *LruRepository[K, T]) getId(entity *T) K {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		panic("ID field not found")
	}
	id, ok := idField.Interface().(K)
	if !ok {
		panic("ID Interface not found")
	}
	return id
}

func (r *LruRepository[K, T]) setId(entity *T, id K) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		panic("ID field not found")
	}
	idField.Set(reflect.ValueOf(id))
}
