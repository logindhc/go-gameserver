package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gameserver/core/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache 是基于 Redis 实现的缓存
type RedisCache[K string | int64, T any] struct {
	client     *redis.Client
	prefix     string
	expiration time.Duration
}

// NewRedisCache 创建一个新的 RedisCache 实例
func NewRedisCache[K string | int64, T any](client *redis.Client, prefix string, expiration time.Duration) *RedisCache[K, T] {
	return &RedisCache[K, T]{client: client, prefix: prefix, expiration: expiration}
}

// Get 从缓存中获取指定 ID 的值
func (r *RedisCache[K, T]) Get(id K) *T {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return nil // 没有找到对应的值
	} else if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Get失败", r.prefix, id))
		return nil
	}

	var entity T
	err = json.Unmarshal([]byte(data), &entity)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Get反序列化失败", r.prefix, id))
		return nil
	}
	return &entity
}

// GetAll 获取缓存中所有的值
func (r *RedisCache[K, T]) GetAll() []*T {
	keys, err := r.client.Keys(context.Background(), fmt.Sprintf("%s:*", r.prefix)).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#GetAll查询失败", r.prefix))
		return nil
	}

	var entitys []*T
	for _, key := range keys {
		data, err := r.client.Get(context.Background(), key).Result()
		if err != nil {
			continue
		}

		var entity T
		err = json.Unmarshal([]byte(data), &entity)
		if err != nil {
			continue
		}
		entitys = append(entitys, &entity)
	}
	return entitys
}

// Put 将值放入缓存中
func (r *RedisCache[K, T]) Put(id K, entity *T) *T {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := json.Marshal(entity)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Put序列化失败", r.prefix, id))
		return nil
	}
	tx := r.client.Set(context.Background(), key, data, r.expiration)
	if tx.Err() != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Put失败", r.prefix, id))
		return nil
	}
	return entity
}

// PutIfAbsent 如果缓存中不存在指定 ID 的值，则将值放入缓存中
func (r *RedisCache[K, T]) PutIfAbsent(id K, entity *T) *T {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := json.Marshal(entity)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v PutIfAbsent序列化失败", r.prefix, id))
		return nil
	}

	// 使用 SETNX 命令实现原子性的 PutIfAbsent
	ok, err := r.client.SetNX(context.Background(), key, data, r.expiration).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v PutIfAbsent失败", r.prefix, id))
		return nil
	}
	if !ok {
		return nil // 已经存在该键
	}
	return entity
}

// Replace 替换缓存中的值
func (r *RedisCache[K, T]) Replace(id K, entity *T) bool {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	data, err := json.Marshal(entity)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Replace序列化失败", r.prefix, id))
		return false
	}

	// 使用 GETSET 命令实现原子性的 Replace
	_, err = r.client.GetSet(context.Background(), key, data).Result()
	if errors.Is(err, redis.Nil) {
		return false // 没有找到对应的键
	} else if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Replace失败", r.prefix, id))
		return false
	}
	return true
}

// Remove 从缓存中移除指定 ID 的值
func (r *RedisCache[K, T]) Remove(id K) {
	key := fmt.Sprintf("%s:%v", r.prefix, id)
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#id:%v Remove失败", r.prefix, id))
		return
	}
}

// Clear 清空缓存，一般不建议使用
func (r *RedisCache[K, T]) Clear() {
	keys, err := r.client.Keys(context.Background(), fmt.Sprintf("%s:*", r.prefix)).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#Clear失败", r.prefix))
		return
	}
	if len(keys) <= 0 {
		return
	}
	err = r.client.Del(context.Background(), keys...).Err()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#Clear失败", r.prefix))
	}
}

// Size 返回缓存中元素的数量
func (r *RedisCache[K, T]) Size() int {
	keys, err := r.client.Keys(context.Background(), fmt.Sprintf("%s:*", r.prefix)).Result()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("%s#size查询失败", r.prefix))
		return 0
	}
	return len(keys)
}
