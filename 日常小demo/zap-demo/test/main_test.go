package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func BenchmarkPrintByFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Printf("%d", 1)
	}
}

func BenchmarkPrintByBuf(b *testing.B) {
	buf := bufio.NewWriter(os.Stdout)
	for i := 0; i < b.N; i++ {
		buf.Write([]byte(strconv.Itoa(1)))
		buf.Flush()
	}
}

/*
BenchmarkPrintByFormat-8	9999         105120 ns/op
BenchmarkPrintByBuf-8       17176         64105 ns/op
*/
