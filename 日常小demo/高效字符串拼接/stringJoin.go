package main

import (
	"bytes"
	"fmt"
	"strings"
)

func methodDirectAdd() {
	s := ""
	for i := 0; i < 100; i++ {
		s += "hello"
	}

}

// Strings.Join
func methodStringJoin() {
	s := make([]string, 100, 100)
	for i := range s {
		s[i] = "hello"
	}
	_ = strings.Join(s[:], "")
}

//先不考虑这个方法
func method3() {
	s := fmt.Sprintf("%s %s", "hello ", "world")
	fmt.Println(s)
}

//strings.Builder
func methodStringsBuilder() {
	builder := strings.Builder{}
	for i := 0; i < 100; i++ {
		builder.WriteString("hello")
	}
	_ = builder.String()
}

//bytesBuffer
func methodBytesBuffer() {
	var buf bytes.Buffer
	for i := 0; i < 100; i++ {
		buf.WriteString("hello")
	}
	_ = buf.String()
}
