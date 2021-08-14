package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"
)

//每帧的格式，没用到
const (
	Max_Row   = 110
	Max_Col   = 150
	Max_Frame = 5773
)

var data [Max_Frame]string

func main() {
	cmd := exec.Command("mode", "con", "cols=150", "lines=55")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	//先加载所有数据
	var wg sync.WaitGroup
	for i := 0; i < Max_Frame; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			filename := fmt.Sprintf("out/ASCII-QQ录屏20210808121625%04d.txt", i)

			f, err := os.Open(filename)
			defer f.Close()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			bs, _ := ioutil.ReadAll(f) //试试每次读取一行
			data[i] = string(bs)
		}(i)
	}
	wg.Wait()
	for i := 0; i < Max_Frame; i++ {
		//now := time.Now()
		//cmd := exec.Command("cls")		这玩意执行时间好长
		//cmd.Run()
		//fmt.Println(time.Since(now).Milliseconds())
		fmt.Println(data[i])
		time.Sleep(38 * time.Millisecond)
	}
}

//初始化终端窗口
func play(index int) {

	//

}
