package repository

import (
	"gorm.io/gorm"
)

type IRepository[K string | int64, T any] interface {
	Get(id K) *T
	GetAll() []*T
	GetOrCreate(id K) *T
	Add(entity *T) *T
	Remove(id K)
	Update(entity *T)
	Flush()
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
}
