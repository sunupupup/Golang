package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//实现一致性哈希算法

//Hash 代表着将[]bye映射为 uint32 类型，正好是 2^32 位数的数
type Hash func([]byte) uint32

//MyMap 包含所有key值
type MyMap struct {
	hash     Hash           //采取依赖注入的方式，允许用于替换成自定义的 Hash 函数,默认为 crc32.ChecksumIEEE 算法。
	replicas int            //虚拟节点倍数 replicas
	keys     []int          //排序好的,哈希环 keys,放所有虚拟节点的的keys
	hashmap  map[int]string //虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
}

//创建一个Map实例
func NewMap(replicas int, fn Hash) *MyMap {
	m := &MyMap{
		replicas: replicas,
		hash:     fn,
		hashmap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//实现真实的添加真实节点/机器的 Add() 方法
func (this *MyMap) Add(keys ...string) {

	//计算key，并且添加replices个虚拟节点进map
	for _, key := range keys {
		for i := 0; i < this.replicas; i++ {
			hash := int(this.hash([]byte(key + strconv.Itoa(i)))) //根据key，再加上虚拟节点索引，计算hash值，strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
			this.keys = append(this.keys, hash)
			this.hashmap[hash] = key //记录虚拟节点（hash环上的int值）和 真实节点（[]byte类型的key）的映射关系
		}
	}
	sort.Ints(this.keys) //最后一步，环上的哈希值排序。

}

//实现选择节点的Get()方法
func (this *MyMap) Get(key string) string {
	//先计算哈希值
	if len(key) == 0 {
		return ""
	}
	hash := this.hash([]byte(key)) //找到第一个大的节点
	index := sort.Search(len(this.keys), func(i int) bool {
		return this.keys[i] >= int(hash)
	})
	//这里的index只是找到一个虚拟节点，还要找到真正的机器节点
	//idx == len(m.keys)时，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	return this.hashmap[this.keys[index%(len(this.keys))]]

}
