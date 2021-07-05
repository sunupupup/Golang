package main

import (
	"bufio"
	"os"
)

func main() {

	buf := bufio.NewWriter(os.Stdout)
	buf.Write([]byte("123123"))
	buf.Flush()
}
