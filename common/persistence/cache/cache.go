package cache

type Cache[K string | int64, V any] interface {
	Get(id K) *V
	GetAll() []*V
	Put(id K, v *V) *V
	PutIfAbsent(id K, v *V) *V
	Replace(id K, v *V) bool
	Remove(id K)
	Clear()
	Size() int
}
