package main

import (
	"fmt"
	"os"
	"runtime/pprof"
)

////////////////////////////////////用法1/////////////////////////////
/*    pkg/profile
func fib(n int) int {
	if n <= 1 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}
func main() {

	defer profile.Start().Stop()

	n := 10
	for i := 1; i <= 5; i++ {
		fmt.Printf("fib(%d)=%d\n", n, fib(n))
		n += 3 * i
	}

}

/////////////////////////////////////用法2/////////////////////////////
*/
/*  mem profiling
func main() {
	f, _ := os.OpenFile("mem.profile", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	for i := 0; i < 500; i++ {
		repeat(generate(1000), 100)
	}

	pprof.Lookup("heap").WriteTo(f, 0)
}

const Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generate(n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.WriteByte(Letters[rand.Intn(len(Letters))])
	}
	return buf.String()
}

func repeat(s string, n int) string {
	var result string
	for i := 0; i < n; i++ {
		result += s
	}

	return result
}
*/

//   CPU  profiling
func main() {
	f, _ := os.OpenFile("cpu.profile", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 9; i++ {
		fmt.Println(fib(i * 5))
	}
}

func fib(n int) int {
	if n <= 1 {
		return 1
	}

	return fib(n-1) + fib(n-2)
}
func fib2(n int) int {
	if n <= 1 {
		return 1
	}

	f1, f2 := 1, 1
	for i := 2; i <= n; i++ {
		f1, f2 = f2, f1+f2
	}

	return f2
}

//////////////////////////////////////用法3/////////////////////////////
//  _   net/http/pprof
