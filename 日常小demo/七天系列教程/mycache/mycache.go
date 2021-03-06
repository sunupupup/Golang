package mycache

import (
	"errors"
	"log"
	"sync"

	pb "github.com/零碎知识/mycache/mycachepb"
	"github.com/零碎知识/mycache/singleflight"
)

//我们思考一下，如果缓存不存在，应从数据源（文件，数据库等）获取数据并添加到缓存中。GeeCache 是否应该支持多种数据源的配置呢？
//不应该，一是数据源的种类太多，没办法一一实现；二是扩展性不好。
//如何从源头获取数据，应该是用户决定的事情，我们就把这件事交给用户好了。
//因此，我们设计了一个回调函数(callback)，在缓存不存在时，调用这个函数，得到源数据。

//设计一个Getter接口，获取key的value
type Getter interface {
	Get(key string) ([]byte, error)
}

//定义函数类型 GetterFunc，并实现 Getter 接口的 Get 方法
type GetterFunc func(key string) ([]byte, error)

//函数类型实现某一个接口，称之为接口型函数,方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数。
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//定义最重要的Group结构体

//一个Group代表一个缓存空间，有自己的name
type Group struct {
	name      string //自己的名字
	getter    Getter //回调函数
	mainCache cache  //并发安全的缓存
	peers     PeerPicker
	//使用 singleflight.Group 确保m每个key只会被去查一次，不会出现同时往db种查同一个key的情况
	loader *singleflight.Group
}

//利用读写锁，实现对 groups的map的读写并发控制
var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//创建一个新的Group实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}

	groups[name] = &g

	return &g
}

func GetGroup(name string) *Group {
	mu.RLock()
	ret := groups[name]
	defer mu.RUnlock()
	return ret
}

//实现Get方法  根据key从缓存中获取value
func (this *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key must be required")
	}
	if v, ok := this.mainCache.get(key); ok {
		return v, nil
	}
	//没从缓存中加载到，尝试去加载内容
	return this.load(key)
}

//利用回调函数加载缓存中没有的内容
//缓存不存在，则调用 load 方法，load 调用 getLocally（分布式场景下会调用 getFromPeer 从其他节点获取）
func (this *Group) load(key string) (value ByteView, err error) {

	//尝试获取远程节点的数据,且通过singleflight模式，确保并发环境下只有一个key会被在db中击中
	//将原来的 load 的逻辑，使用 g.loader.Do 包裹起来即可，这样确保了并发场景下针对相同的 key，load 过程只会调用一次。
	viewi, err := this.loader.Do(key, func() (interface{}, error) {
		if this.peers != nil {

			if peer, ok := this.peers.PickPeer(key); ok {
				if value, err = this.getFromPeer(peer, key); err != nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return this.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (this *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {

	req := &pb.Request{
		Group: this.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil

}

//getLocally 调用用户回调函数 g.getter.Get() 获取源数据，并且将源数据添加到缓存 mainCache 中（通过 populateCache 方法）
func (this *Group) getLocally(key string) (ByteView, error) {

	bytes, err := this.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	ret := ByteView{bytes}
	//要将这部分内容添加到缓存中
	this.populateCache(key, ret)
	return ret, nil
}

func (this *Group) populateCache(key string, value ByteView) {
	this.mainCache.add(key, value)
}

//注册peers节点
func (this *Group) RegisterPeers(peers PeerPicker) {
	if this.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	this.peers = peers

}
