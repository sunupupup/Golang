package main

import "testing"

/*
测试结果：
>> go test -bench=. -benchmem	       	-benchmem查看内存分配 和 每次执行需要alloc分配多少次内存

goos: windows
goarch: amd64
pkg: stringJoin
BenchmarkMethodDirectAdd-8                         67839             27273 ns/op           26384 B/op         99 allocs/op
BenchmarkMethodStringJoin-8                       479966              2446 ns/op             512 B/op          1 allocs/op
BenchmarkMethodStringsBuilder-8                   999849              1061 ns/op            1016 B/op          7 allocs/op
BenchmarkMethodBytesBuffer-8                      600012              2045 ns/op            1584 B/op          5 allocs/op
BenchmarkMethodByteSliceAppendNotAllocated-8     1000000              1168 ns/op            1528 B/op          8 allocs/op
BenchmarkMethodByteSliceAppend-8                 1963845               696 ns/op             512 B/op          1 allocs/op
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

func BenchmarkMethodByteSliceAppendNotAllocated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodByteSliceAppendNotAllocated()
	}
}

//append方式
func BenchmarkMethodByteSliceAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		methodByteSliceAppend()
	}
}
