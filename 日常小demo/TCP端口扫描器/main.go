package main

import (
	"fmt"
	"net"
	"sort"
)

//goroutine池/worker池 并发tcp扫描器
func main() {
	ports := make(chan int, 100)
	result := make(chan int)
	var openPorts []int
	//var closePorts []int
	//创建worker池子，控制goroutine数量
	for i := 0; i < cap(ports); i++ { //相当于开启100个worker，等待传入数据干活
		go worker(ports, result)
	}

	//创建任务
	go func() {
		for i := 1; i < 1024; i++ {
			ports <- i
		}
	}()

	//接受结果
	for i := 1; i < 1024; i++ {
		port := <-result //等待1024次，就不需要waitGroup了
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(result)
	fmt.Println("执行完毕")
	sort.Ints(openPorts)
	fmt.Println(openPorts)
	//结果[135 445]
	//135端口主要用于使用RPC（Remote Procedure Call，远程过程调用）协议并提供DCOM（分布式组件对象模型）服务
	//445端口有了它我们可以在局域网中轻松访问各种共享文件夹或共享打印机
}

func worker(ports chan int, result chan int) {
	for p := range ports {
		addr := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			result <- 0
			continue
		}
		conn.Close()
		result <- p
	}
}

//并发扫描端口
/*
func main() {
	//端口有6w多个，暂时只扫描100多个
	var wg sync.WaitGroup
	start := time.Now()
	for i := 6300; i < 6400; i++ { //127.0.0.1:6379端口打开了  mysql的端口
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			addr := fmt.Sprintf("127.0.0.1:%d", j)
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				//要么端口关闭了，要么端口被防火墙过滤了
				fmt.Printf("%s端口关闭了 \n", addr)
				return
			}
			conn.Close()
			fmt.Printf("%s端口打开了 \n", addr)
		}(i) //必须要传i进去，不然所有的goroutine使用的i都是120

	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("扫描完毕")
	fmt.Println("耗时", elapsed)

}
*/

//非并发的，单线程，很慢
/*
func main() {
	//端口有6w多个，暂时只扫描ip地址的前100多个
	for i := 21; i < 120; i++ {
		addr := fmt.Sprintf("127.0.0.1:%d", i)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			//要么端口关闭了，要么端口被防火墙过滤了
			fmt.Printf("%s端口关闭了 \n", addr)
			continue
		}
		conn.Close()
		fmt.Printf("%s端口打开了 \n", addr)
	}
}
*/
