package mycache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

//测试回调函数
func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	target := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, target) {
		t.Error("callback failed")
	}
}

//测试样例
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {

	loadCount := make(map[string]int)

	//获取一个新的缓存组
	group := NewGroup("students", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			//如果db中存在
			if v, ok := db[key]; ok {
				//查看这里被调用了多少次
				if _, ok := loadCount[key]; !ok {
					loadCount[key] = 0
				}
				loadCount[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s 不存在", key)
		},
	))

	for k, v := range db {
		//尝试获取，应该是一遍成功，一遍失败
		if get, err := group.Get(k); err != nil || get.String() != v {
			t.Fatal("failed to get value of Tom")
		}

		if _, err := group.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatal("failed to get value of Tom")
		}
	}

	if get, err := group.Get("unknown"); err == nil {
		t.Fatal("key:unknown 应该不存在，可是获取到了value:", get)
	}

}
