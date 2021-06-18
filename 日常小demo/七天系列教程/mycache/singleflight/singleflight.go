package singleflight

import "sync"

//call 代表正在进行中，或者已经结束的请求，使用sync.WaitGroup 锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//Group是singleflight的主要数据结构，管理不同的key的请求（call）
type Group struct {
	mu sync.Mutex //保护 m
	m  map[string]*call
}

//实现Do方法
func (this *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {

	this.mu.Lock()
	if this.m == nil {
		this.m = make(map[string]*call)
	}
	//获取到call
	if c, ok := this.m[key]; ok {
		this.mu.Unlock()
		c.wg.Wait() //如果这个key对应的call已经发起，就等待上一个call结束，并直接返回结果
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	this.m[key] = c
	this.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	this.mu.Lock()
	delete(this.m, key)		//查完了就没必要在map里放着了，已经放到缓存中了
	this.mu.Unlock()

	return c.val, c.err

}
