package channel

type SafeMap struct {
	data     map[string]interface{}
	opsChan  chan func() // 用于接收操作的 channel
	doneChan chan bool   // 用于关闭后台 Goroutine
}

// NewSafeMap 初始化并返回一个使用 channel 实现的 SafeMap
func NewSafeMap() *SafeMap {
	sm := &SafeMap{
		data:     make(map[string]interface{}),
		opsChan:  make(chan func()),
		doneChan: make(chan bool),
	}

	// 启动后台 Goroutine 处理所有操作
	go sm.run()
	return sm
}

// 后台 Goroutine 处理所有操作
func (sm *SafeMap) run() {
	for {
		select {
		case op := <-sm.opsChan:
			op() // 执行传入的操作
		case <-sm.doneChan:
			return // 关闭 Goroutine
		}
	}
}

// Close 关闭 SafeMap 的后台 Goroutine
func (sm *SafeMap) Close() {
	close(sm.doneChan)
}

// Set 设置一个键值对
func (sm *SafeMap) Set(key string, value interface{}) {
	// 通过 opsChan 发送写入操作
	sm.opsChan <- func() {
		sm.data[key] = value
	}
}

// Get 获取指定键的值
func (sm *SafeMap) Get(key string) (interface{}, bool) {
	resultChan := make(chan interface{}) // 用于接收结果
	sm.opsChan <- func() {
		value, exists := sm.data[key]
		if exists {
			resultChan <- value
		} else {
			resultChan <- nil
		}
		close(resultChan)
	}
	value := <-resultChan
	if value == nil {
		return nil, false
	}
	return value, true
}

// Delete 删除一个键
func (sm *SafeMap) Delete(key string) {
	// 通过 opsChan 发送删除操作
	sm.opsChan <- func() {
		delete(sm.data, key)
	}
}
