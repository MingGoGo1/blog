package global

import "sync"

// Cache 全局缓存变量
var Cache map[string]any

// cacheMutex 缓存操作的互斥锁
var cacheMutex sync.RWMutex

// InitCache 初始化缓存
func InitCache() {
	Cache = make(map[string]any)
}

// SetCache 设置缓存
func SetCache(key string, value any) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	Cache[key] = value
}

// GetCache 获取缓存
func GetCache(key string) (any, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	value, exists := Cache[key]
	return value, exists
}

// DeleteCache 删除缓存
func DeleteCache(key string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	delete(Cache, key)
}
