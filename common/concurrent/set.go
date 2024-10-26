package concurrent

import (
	"sync"
)

// ConcurrentSet 是一个线程安全的集合，可以添加和删除元素。
type ConcurrentSet[T comparable] struct {
	set sync.Map
}

// NewConcurrentSet 创建一个新的 ConcurrentSet 实例。
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{sync.Map{}}
}

// Add 向集合中添加一个元素。
func (s *ConcurrentSet[T]) Add(item T) {
	s.set.Store(item, struct{}{})
}

// Remove 从集合中移除一个元素。
func (s *ConcurrentSet[T]) Remove(item T) {
	s.set.Delete(item)
}

// Has 检查集合中是否存在一个元素。
func (s *ConcurrentSet[T]) Has(item T) bool {
	_, exists := s.set.Load(item)
	return exists
}

// Clear 清空集合中的所有元素。
func (s *ConcurrentSet[T]) Clear() {
	s.set.Clear()
}

// Size 返回集合中元素的数量。
func (s *ConcurrentSet[T]) Size() int {
	siz := 0
	s.set.Range(func(k, v any) bool {
		siz++
		return true
	})
	return siz
}

// All 返回一个包含集合中所有元素的切片。
func (s *ConcurrentSet[T]) All() []T {
	// 获取所有键
	var keys []T
	s.set.Range(func(k, v any) bool {
		keys = append(keys, k.(T))
		return true
	})
	return keys
}
