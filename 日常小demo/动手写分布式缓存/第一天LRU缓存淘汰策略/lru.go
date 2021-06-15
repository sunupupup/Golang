//实现 LRU 缓存淘汰算法
//1. FOFO 先进先出
//2. LRU  Least Recently   Used   最近最少使用
//        LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了，则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。
//3. LFU  Least Frequently Used   最少使用，使用频率最低，
package lru

import (
	"container/list"
)

//LRU，认为最近被访问的数据，以后被访问的概率也很大
//使用一个队列维护缓存，实现的话用map和双向链表存放
//此时还是并发不安全的
type Cache struct {
	maxBytes  int64 //maxBytes 是允许使用的最大内存，超过最大内存就要淘汰缓存了
	nbytes    int64 //nbytes 是当前已使用的内存，只记录 key 和 value 的开销
	ll        *list.List
	cache     map[string]*list.Element      //键是字符串，值是双向链表中对应节点的指针。定义Element中的Value是entry类型
	OnEvicted func(key string, calue Value) //OnEvicted 是某条记录被移除时的回调函数，可以为 nil。
}

//键值对 entry 是双向链表节点的数据类型，
//在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value //完全可以设置为一个string，方便扩展，就用Value了，只要实现Len方法就好
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

//新建一个cache对象
func NewCache(maxbytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxbytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查找功能
//如果存在，就移到队头
func (this *Cache) Get(key string) (value Value, ok bool) {
	//通过map查找
	if ele, ok := this.cache[key]; ok {
		this.ll.MoveToFront(ele) //将element移到链表头
		data := ele.Value.(*entry)
		return data.value, true
	}
	return nil, false
}

//删除，也就是淘汰最老的缓存，也就是队尾
func (this *Cache) RemoveOldest() {
	ele := this.ll.Back()
	//还要从map中删除
	if ele != nil {
		kv := ele.Value.(*entry)
		delete(this.cache, kv.key)
		//更新Cache中的nbytes字段
		this.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//如果回调函数不空，还要调用
		if this.OnEvicted != nil {
			this.OnEvicted(kv.key, kv.value)
		}
	}
}

//新增/修改缓存
func (this *Cache) Add(key string, value Value) {
	//往链表头加
	//先判断有没有
	if ele, ok := this.cache[key]; ok {
		//如果已经有了,更新数据，并放到队头
		this.ll.MoveToFront(ele)
		//更新nbytes大小  现在的长度 - 原来的长度
		this.nbytes += int64(value.Len()) - int64(ele.Value.(*entry).value.Len())
		//更新数据
		ele.Value.(*entry).value = value
	} else {
		//先添加进去，如果超出内存了，后续再淘汰
		ele = this.ll.PushFront(&entry{key, value})
		this.cache[key] = ele
		this.nbytes += int64(len(key)) + int64(value.Len())
	}

	//如果超出缓存上限，进行淘汰
	for this.nbytes > this.maxBytes {
		//fmt.Println("淘汰内存了")
		this.RemoveOldest()
	}

	//temp, ok := this.cache[key]
	//return temp.Value.(*entry).value, ok

}

// Len 统计多少条内存
func (c *Cache) Len() int {
	return c.ll.Len()
}

/*
type String string

func (this String) Len() int {
	return len(this)
}
func main() {

	lru := NewCache(int64(1024), nil)
	fmt.Println(lru == nil)
	lru.Add("key1", String("1234"))
	fmt.Println(lru.nbytes)
	v, _ := lru.Get("key1")
	fmt.Println(v.(String))
	//if ok {
	//	fmt.Println(v.(String))
	//}
	//fmt.Println(s, ok)

		if !ok {
			fmt.Println("不ok")
		}
		fmt.Println(lru.Len())
		fmt.Println(111)
		fmt.Println(s)
		if v, ok := lru.Get("key1"); ok {
			fmt.Println("get the key,and the value is :", v.(String))
		} else {
			fmt.Println("111")
		}

}

*/
