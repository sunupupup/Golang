//基准测试
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
	"time"
)

//基准结果对应的结构体
type BenchmarkResult struct {
	N         int           // 迭代次数
	T         time.Duration // 基准测试花费的时间
	Bytes     int64         // 一次迭代处理的字节数
	MemAllocs uint64        // 总的分配内存的次数
	MemBytes  uint64        // 总的分配内存的字节数

	// Extra records additional metrics reported by ReportMetric.
	Extra map[string]float64
}

//执行命令  go test -benchmem -bench .  输出
//BenchmarkHello-8   	10523794	       147 ns/op	       5 B/op	       1 allocs/op

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hello")
	}
}

//如果在运行前需要一些配置，那么还需要重置定时器、
func BenchmarkHello2(b *testing.B) {
	time.Sleep(time.Second)
	b.ResetTimer() //重新计时
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("hello")
	}
}

//使用 RunParallel 测试并发性能
func BenchmarkParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	//The number of goroutines defaults to GOMAXPROCS
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}
