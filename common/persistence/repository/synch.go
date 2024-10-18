package repository

import (
	"gameserver/common/persistence/buffer"
	"gorm.io/gorm"
	"reflect"
	"sync"
)

type SynchRepository[K string | int64, T any] struct {
	db     *gorm.DB
	buffer *buffer.SyncBuffer[K, T]
	sync.Mutex
}

func NewSynchRepository[K string | int64, T any](db *gorm.DB) *SynchRepository[K, T] {
	r := &SynchRepository[K, T]{
		db:     db,
		buffer: buffer.NewSyncBuffer[K, T](db),
	}
	return r
}

func (r *SynchRepository[K, T]) Get(id K) *T {
	var entity T
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	tx := r.db.Where("id = ?", id).Find(&entity)
	if tx.RowsAffected == 0 {
		return nil
	}
	return &entity
}

func (r *SynchRepository[K, T]) GetAll() []*T {
	var entities []*T
	r.db.Find(&entities)
	return entities
}

func (r *SynchRepository[K, T]) GetOrCreate(id K) *T {
	entity := r.Get(id)
	if entity == nil {
		entity = new(T)
		r.setId(entity, id)
		return r.Add(entity)
	}
	return entity
}

func (r *SynchRepository[K, T]) Add(entity *T) *T {
	if entity == nil {
		return nil
	}
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.buffer.Add(entity)
	return entity
}

func (r *SynchRepository[K, T]) Remove(id K) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.buffer.Remove(id)
}

func (r *SynchRepository[K, T]) Update(entity *T) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.buffer.Update(entity)
}
func (r *SynchRepository[K, T]) Flush() {
}

func (r *SynchRepository[K, T]) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return r.db.Where(query, args...)
}

func (r *SynchRepository[K, T]) getId(entity *T) K {
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

func (r *SynchRepository[K, T]) setId(entity *T, id K) {
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
