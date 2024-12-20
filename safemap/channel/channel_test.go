package channel

import (
	"sync"
	"testing"
)

func TestSafeMapConcurrency(t *testing.T) {
	sm := NewSafeMap()
	defer sm.Close() // 确保测试完成后关闭 SafeMap

	wg := sync.WaitGroup{}

	// 并发写
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(string(rune(i)), i)
		}(i)
	}

	// 并发读
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Get(string(rune(i)))
		}(i)
	}

	wg.Wait()

	// 验证所有键是否写入成功
	for i := 0; i < 100; i++ {
		_, exists := sm.Get(string(rune(i)))
		if !exists {
			t.Errorf("Key %v missing after concurrent writes", i)
		}
	}
}
