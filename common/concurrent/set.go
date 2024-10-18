package concurrent

import (
	"sync"
)

// ConcurrentSet 是一个线程安全的集合，可以添加和删除元素。
type ConcurrentSet[T comparable] struct {
	lock  sync.RWMutex
	items map[T]struct{}
}

// NewConcurrentSet 创建一个新的 ConcurrentSet 实例。
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{
		lock:  sync.RWMutex{},
		items: make(map[T]struct{}),
	}
}

// Add 向集合中添加一个元素。
func (set *ConcurrentSet[T]) Add(item T) {
	set.lock.Lock()
	defer set.lock.Unlock()
	set.items[item] = struct{}{}
}

// Remove 从集合中移除一个元素。
func (set *ConcurrentSet[T]) Remove(item T) {
	set.lock.Lock()
	defer set.lock.Unlock()
	delete(set.items, item)
}

// Has 检查集合中是否存在一个元素。
func (set *ConcurrentSet[T]) Has(item T) bool {
	set.lock.RLock()
	defer set.lock.RUnlock()
	_, exists := set.items[item]
	return exists
}

// Clear 清空集合中的所有元素。
func (set *ConcurrentSet[T]) Clear() {
	set.lock.Lock()
	defer set.lock.Unlock()
	set.items = make(map[T]struct{})
}

// Size 返回集合中元素的数量。
func (set *ConcurrentSet[T]) Size() int {
	set.lock.RLock()
	defer set.lock.RUnlock()
	return len(set.items)
}

// All 返回一个包含集合中所有元素的切片。
func (set *ConcurrentSet[T]) All() []T {
	set.lock.RLock()
	defer set.lock.RUnlock()
	// 获取所有键
	var keys []T
	for key := range set.items {
		keys = append(keys, key)
	}
	return keys
}
