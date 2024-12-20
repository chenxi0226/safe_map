package rwmutex

import (
	"sync"
	"testing"
)

func TestSafeMapConcurrency(t *testing.T) {
	sm := NewSafeMap()
	wg := sync.WaitGroup{}

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(string(rune(i)), i)
		}(i)
	}

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Get(string(rune(i)))
		}(i)
	}

	wg.Wait()

	// 验证所有键是否写入
	for i := 0; i < 100; i++ {
		_, exists := sm.Get(string(rune(i)))
		if !exists {
			t.Errorf("Key %v missing after concurrent writes", i)
		}
	}
}
