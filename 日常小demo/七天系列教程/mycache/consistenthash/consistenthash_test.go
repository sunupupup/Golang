package consistenthash

import (
	"fmt"
	"testing"
)

func TestHashing(t *testing.T) {

	//先构建一个MyMap对象
	m := NewMap(160, nil)
	m.Add("Node1", "Node2", "Node3")
	fmt.Println(len(m.keys))
	testcase := map[string]string{
		"Node10":   "Node1",
		"Node1159": "Node1", //Node10 --  Node1159 都是 Node1的虚拟子节点
		"Node20":   "Node2",
		"Node30":   "Node3",
	}
	for k, v := range testcase {
		if m.Get(k) != v {
			t.Errorf("key is %s,target is %s,got %s", k, m.Get(k), v)
		}
	}

	//fmt.Println(m.Get("Node40"))	//Node2

	m.Add("Node4")
	if m.Get("Node40") != "Node4" {
		t.Error("增加节点失败")
	}

}
