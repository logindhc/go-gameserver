package concurrent

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func strfnv32[T comparable](key T) uint32 {
	return fnv32(fmt.Sprintf("%v", key))
}

// ConcurrentSet 是一个线程安全的集合，可以添加和删除元素。
type ConcurrentSet[T comparable] struct {
	cmap.ConcurrentMap[T, struct{}]
}

// NewConcurrentSet 创建一个新的 ConcurrentSet 实例。
func NewConcurrentSet[T comparable]() *ConcurrentSet[T] {
	return &ConcurrentSet[T]{cmap.NewWithCustomShardingFunction[T, struct{}](strfnv32)}
}

// Add 向集合中添加一个元素。
func (s *ConcurrentSet[T]) Add(item T) {
	s.Set(item, struct{}{})
}

// Remove 从集合中移除一个元素。
func (s *ConcurrentSet[T]) Remove(item T) {
	s.Remove(item)
}

// Has 检查集合中是否存在一个元素。
func (s *ConcurrentSet[T]) Has(item T) bool {
	exists := s.Has(item)
	return exists
}

// Clear 清空集合中的所有元素。
func (s *ConcurrentSet[T]) Clear() {
	s.Clear()
}

// Size 返回集合中元素的数量。
func (s *ConcurrentSet[T]) Size() int {
	count := s.Count()
	return count
}

// All 返回一个包含集合中所有元素的切片。
func (s *ConcurrentSet[T]) All() []T {
	return s.Keys()
}
