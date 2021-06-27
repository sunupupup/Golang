package main

import "testing"

/*
测试结果：
>> go test -bench=. -benchmem	       	-benchmem查看内存分配 和 每次执行需要alloc分配多少次内存

goos: windows
goarch: amd64
pkg: stringJoin
BenchmarkMethodDirectAdd-8                 68960             16366 ns/op           26384 B/op         99 allocs/op
BenchmarkMethodStringJoin-8               600000              2037 ns/op             512 B/op          1 allocs/op
BenchmarkMethodStringsBuilder-8          1015584              1274 ns/op            1016 B/op          7 allocs/op
BenchmarkMethodBytesBuffer-8              571411              2341 ns/op            1584 B/op          5 allocs/op
*/

func BenchmarkMethodDirectAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodDirectAdd()
	}
}

func BenchmarkMethodStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodStringJoin()
	}
}

func BenchmarkMethodStringsBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodStringsBuilder()
	}
}

func BenchmarkMethodBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodBytesBuffer()
	}
}
