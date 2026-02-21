package utils

import (
	"sync"
	"time"
)

// cacheItem 缓存项，包含数据和过期时间
type cacheItem struct {
	data      interface{}
	expiresAt time.Time
}

// StatsCache 内存缓存
type StatsCache struct {
	store sync.Map
}

var GlobalCache = &StatsCache{}

// Set 设置缓存，有效期 1 小时
func (c *StatsCache) Set(key string, data interface{}) {
	c.store.Store(key, cacheItem{
		data:      data,
		expiresAt: time.Now().Add(1 * time.Hour), // 严格缓存 1 小时
	})
}

// Get 获取缓存，如果过期则返回 nil
func (c *StatsCache) Get(key string) interface{} {
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	item := val.(cacheItem)
	if time.Now().After(item.expiresAt) {
		c.store.Delete(key) // 过期删除
		return nil
	}

	return item.data
}

// GenerateKey 生成带 user_id 的缓存 key
func GenerateKey(userID string, prefix string, suffix string) string {
	return userID + ":" + prefix + ":" + suffix
}
