package main

/*
结果： switch比map快多了，但是观感就差多了，而且需要大量分支，代码复杂
BenchmarkMap-8          23066143                49.8 ns/op
BenchmarkSwitch-8       1000000000               0.738 ns/op
*/

import "testing"

const (
	State1 = iota + 1  //按照uber的规范，最好是 itoa+1 开始
	State2
	State3
	State4
)

var m = map[int]string{
	State1: "状态1",
	State2: "状态2",
	State3: "状态3",
	State4: "状态4",
}

func testMap(state int) string {
	return m[state]
}

func testSwitch(state int) string {
	var s string
	switch state {
	case State1:
		s = "状态1"
	case State2:
		s = "状态2"
	case State3:
		s = "状态3"
	case State4:
		s = "状态4"
	}
	return s
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testMap(0)
		testMap(1)
		testMap(2)
		testMap(3)
	}
}

func BenchmarkSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSwitch(0)
		testSwitch(1)
		testSwitch(2)
		testSwitch(3)
	}
}
