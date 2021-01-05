package singleflight

import (
	"sync"
)

// call is an in-flight or completed Do call
// 代表正在进行中，或已经结束的请求
type call struct {
	wg  sync.WaitGroup // 并发协程之间不需要消息传递，非常适合 sync.WaitGroup
	val interface{}
	err error
}

// Group 管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call // lazily initialized
}

// Do 针对相同的key，无论Do被调用多少次，函数fn都只会被调用一次，等待 fn 调用结束了，返回返回值或错误
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 有key对应的请求（call）
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()         // 如果请求正在进行中，则等待,等待其他的某个goroutine返回，用wg即可不同通信告知结束
		return c.val, c.err // 请求结束，返回结果
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c // 添加到 g.m，表明 key 已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新 g.m，同一请求时间内的请求，只会响应一次，之后有请求再响应
	g.mu.Unlock()

	return c.val, c.err
}
