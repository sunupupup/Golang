package mycache

import (
	"sync"

	lru "github.com/零碎知识/mycache/lru"
)

//两个目标：
//1.sync.Mutex的使用，并实现之前的lru缓存的并发控制
//2.实现MyCache核心数据结构 Group ，缓存不存在时，调用回调函数获取源数据

//添加并发特性
type cache struct {
	m          sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

//并发安全 添加缓存
func (this *cache) add(key string, value ByteView) {
	this.m.Lock()
	defer this.m.Unlock()
	if this.lru == nil {
		this.lru = lru.NewCache(this.cacheBytes, nil)
	}
	this.lru.Add(key, value)
}

//并发安全获取缓存
func (this *cache) get(key string) (value ByteView, ok bool) {
	this.m.Lock()
	defer this.m.Unlock()
	if this.lru == nil {
		this.lru = lru.NewCache(this.cacheBytes, nil)
	}

	if v, ok := this.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

//设计主体结构：Group，负责与用户的交互，并且控制缓存值存储和获取的流程
/*
                        是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/
//我们将在 geecache.go 中实现主体结构 Group，那么 GeeCache 的代码结构的雏形已经形成了。
/*
geecache/
    |--lru/
        |--lru.go  // lru 缓存淘汰策略
    |--byteview.go // 缓存值的抽象与封装
    |--cache.go    // 并发控制
    |--geecache.go // 负责与外部交互，控制缓存存储和获取的主流程
*/
