package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	test1()
	test2()
}

func test1() {
	var s = [64 * 1024 * 1024]int{}
	fmt.Println(unsafe.Sizeof(int(1)))
	for j := 0; j < 5; j++ {
		now := time.Now()
		for i := 0; i < len(s); i++ {
			s[i] = 1
		}
		fmt.Printf("耗时:%d毫秒", time.Now().Sub(now).Milliseconds())
		now = time.Now()
		for i := 0; i < len(s); i += 8 {
			s[i] = 1
		}
		fmt.Printf("耗时:%d毫秒\n", time.Now().Sub(now).Milliseconds())
	}
}

/*
输出结果：	可以看到 基本没有差多少时间  问题：第一次怎么花了这么久？？
8
耗时:421毫秒耗时:95毫秒
耗时:100毫秒耗时:86毫秒
耗时:100毫秒耗时:94毫秒
耗时:112毫秒耗时:93毫秒
耗时:111毫秒耗时:96毫秒
*/

func test2() {
	var s = [64 * 1024 * 1024]int{}
	fmt.Println(unsafe.Sizeof(int(1)))
	for j := 0; j < 5; j++ {
		now := time.Now()
		for i := 0; i < len(s); i++ {
			s[i] = 1
		}
		fmt.Printf("耗时:%d毫秒", time.Now().Sub(now).Milliseconds())
		now = time.Now()
		for i := 0; i < len(s); i += 16 { //
			s[i] = 1
		}
		fmt.Printf("耗时:%d毫秒\n", time.Now().Sub(now).Milliseconds())
	}
}

/*  这时消耗减少了一半
耗时:325毫秒耗时:55毫秒
耗时:105毫秒耗时:55毫秒
耗时:100毫秒耗时:53毫秒
耗时:112毫秒耗时:60毫秒
耗时:96毫秒耗时:52毫秒
*/

//解释：这与CPU cacheline有关
/*
在本例中 int是8个字节，cacheline是64个字节
所以取64B中的任何一个数据都要去交互一行，所以 +1  和 +8 几乎没有区别
因为他们中间进行缓存交互的这一部分数据是一摸一样的，都要把所有数据缓存到cacheline中

而+1  +16 减少了一半的交互行数，运行耗时就少了一半

所以 数据的交互是很耗时的，数据计算部分是很快的
耗时主要是在缓存这部分，而不是CPU计算这部分


*/
