package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type LRUCache[K string | int64, V any] struct {
	*cache.Cache
}

func NewLRUCache[K string | int64, V any](expiration time.Duration, cleanupInterval time.Duration) *LRUCache[K, V] {
	return &LRUCache[K, V]{cache.New(expiration, cleanupInterval)}
}

func (l *LRUCache[K, V]) Get(id K) *V {
	v, ok := l.Cache.Get(string(id))
	if !ok {
		return nil
	}
	return v.(*V)
}

func (l *LRUCache[K, V]) GetAll() []*V {
	var result []*V
	for _, v := range l.Cache.Items() {
		result = append(result, v.Object.(*V))
	}
	return result
}

func (l *LRUCache[K, V]) Put(id K, value *V) *V {
	l.Cache.SetDefault(string(id), value)
	return value
}

func (l *LRUCache[K, V]) PutIfAbsent(id K, value *V) *V {
	err := l.Cache.Add(string(id), value, cache.DefaultExpiration)
	if err != nil {
		return l.Get(id)
	}
	return value
}

func (l *LRUCache[K, V]) Replace(id K, value *V) bool {
	err := l.Cache.Replace(string(id), value, cache.DefaultExpiration)
	if err != nil {
		return false
	}
	return true
}

func (l *LRUCache[K, V]) Remove(key K) {
	l.Cache.Delete(string(key))
}

func (l *LRUCache[K, V]) Clear() {
	l.Cache.Flush()
}

func (l *LRUCache[K, V]) Size() int {
	return l.Cache.ItemCount()
}
