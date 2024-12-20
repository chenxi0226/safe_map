package rwmutex

import (
	"sync"
)

// SafeMap 是使用 RWMutex 实现的线程安全 Map
type SafeMap struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewSafeMap 初始化并返回一个 SafeMap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]interface{}),
	}
}

// Set 添加或更新一个键值对
func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// Get 获取一个键的值
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, exists := sm.data[key]
	return value, exists
}

// Delete 删除一个键
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

// Keys 返回所有键
func (sm *SafeMap) Keys() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	keys := make([]string, 0, len(sm.data))
	for key := range sm.data {
		keys = append(keys, key)
	}
	return keys
}
