package main

import (
	"strconv"
)

func main() {

	bs := []byte{}
	bs = strconv.AppendInt(bs, 999, 10) //表示将999转为10进制的字符串，并且直接append到[]byte中

	//看源码发现 0-99的数字，是通过硬编码的形式直接解析成string的
}
