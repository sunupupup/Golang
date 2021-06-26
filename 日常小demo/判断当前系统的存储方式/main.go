package main

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

//用Go判断，当前操作系统是大端存储还是小端存储
func main() {

	fn0()

	//通过位运算即可实现
	//大小端转换，这涉及到网络传输、文件存储，要保证读到的数据正确
	//利用标准库encoding/binary
	BigEndianAndLittleEndianByLibrary()
}

func fn0() {
	var a int64 = 1

	//低32位
	a1p := (*int32)(unsafe.Pointer(&a))
	fmt.Println(*a1p)

	//高32位
	a2p := (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + uintptr(unsafe.Sizeof(int32(1)))))
	fmt.Println(*a2p)

	//打印结果 1  0    明显是 低位字节放在低地址端，所以是小端存储
	//大端存储：高低高低，也就是高位字节排放在内存低地址端，高地址端存在低位字节；
	//小端存储：高高低低，也就是高位字节排放在内存的高地址端，低位字节排放在内存的低地址端；
	/*
		0x1A2B3C4D在大端与小端的表现形式

		大端		0x0400		0x4001		0x4002		0x4003				小端 	0x0400		0x4001		0x4002		0x4003

					 1A			  2B		  3C		 4D							 4D			  3C		  2B 		 1A
	*/
}

// uint32 位的大小端转化	位运算
func SwapEndianUint32(val uint32) uint32 {
	return (val&0xff000000)>>24 | (val&0x00ff0000)>>8 |
		(val&0x0000ff00)<<8 | (val&0x000000ff)<<24
}

//官方库
func BigEndianAndLittleEndianByLibrary() {
	var value uint32 = 10
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, value)
	fmt.Println("转化为大端之后的byte数组:", bs)
	fmt.Println("使用大端转化输出之后的结果", binary.BigEndian.Uint32(bs))
	little := binary.LittleEndian.Uint32(bs) //将大端的数组转化为小端的值
	fmt.Println("大端字节序使用小端输出结果：", little)
}
