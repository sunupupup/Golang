package lru

import (
	"fmt"
	"strconv"
	"testing"
)

type String string

func (this String) Len() int {
	return len(this)
}

func TestAddAndGet(t *testing.T) {

	lru := NewCache(int64(1024), nil)
	lru.Add("key1", String("123456"))
	if v, ok := lru.Get("key1"); !ok || v.(String) != String("123456") {
		t.Fatal("test error")
	}

	//不存在的值
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("test error")
	}

}

func TestAddToMax(t *testing.T) {
	//超出内存上限测试
	lru := NewCache(int64(1024), nil)
	for i := 0; i < 100; i++ {
		lru.Add(fmt.Sprintf("key%d", i), String("testtest"+strconv.Itoa(i)))
		if lru.nbytes > lru.maxBytes {
			t.Fatal("缓存内存超了")
		}
	}

}
