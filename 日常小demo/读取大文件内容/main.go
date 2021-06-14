package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

/*
//案例说明：
//使用Golang，25秒读取16GB日志文件，从中提取特定内容   参考链接：https://mp.weixin.qq.com/s/v7YMoTYaYn5jHULEN2pDWg
打开文件后有两个选择：
1. 逐行读取文件，有助于内存紧张，但需要更多的IO
2. 一次将内容全部读取并处理，需要非常大的内存，但总体时间会减少很多
3. bufio.NewReader()将文件分块加载，并且利用sync.Pool 复用4k大小的buf对象; 需要注意的就是要读完整的一行
*/
var wg sync.WaitGroup

func main() {

	bufpool := sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 4*1024)
			return buf
		},
	}

	f, err := os.Open("bigFile.txt")
	if err != nil {
		log.Fatalln(err.Error())
	}

	reader := bufio.NewReader(f)

	for {
		buf := bufpool.Get().([]byte)
		n, err := reader.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err != nil {
				log.Panicln(err.Error())
			}
			if err == io.EOF {
				log.Println(err.Error())
				break
			}
		}
		//读完这一行
		completeLine, err := reader.ReadBytes('\n') //把这一行读完
		if err != io.EOF {
			buf = append(buf, completeLine...)
		}

		//开go程处理这部分数据
		wg.Add(1)
		go processBuf(buf, &bufpool)

	}

}

func processBuf(buf []byte, bufpool *sync.Pool) {
	defer wg.Done()

	var wg2 sync.WaitGroup
	//处理这个buf内容
	//优化：这里也可以复用 string 对象、[]string 对象

	//读取buf的一行一行内容，checkLine检查每一行内容
	lines := string(buf)
	(*bufpool).Put(buf) //复用buf对象
	slicelines := strings.Split(lines, "\n")
	for i := 0; i < len(slicelines); i++ {
		wg2.Add(1)
		go checkLine(slicelines[i], wg2)
	}
	wg2.Wait()
}

func checkLine(line string, wg2 sync.WaitGroup) {
	defer wg2.Done()
	//TODO:
	//检查这一行的数据是否符合标准

}

//bufio用法   但要注意的是，读了4k，并没有到一行的末尾，就可能导致日志一行没读完，需要进一步的操作,main中补充
func test() {
	f, err := os.Open("bigFile.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bufio.NewReader(f)
	//循环读取文件内容
	for {

		buf := make([]byte, 1024*4) //每次读取4k的内容
		n, err := reader.Read(buf)
		if n == 0 {
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if err == io.EOF {
				break
			}
		}
		fmt.Println(buf) //4k的文件内容

	}
}

//以下是别人完整的代码
//我们将根据命令行提供的时间戳提取日志。
func main2() {

	s := time.Now()
	args := os.Args[1:]
	if len(args) != 6 { // for format  LogExtractor.exe -f "From Time" -t "To Time" -i "Log file directory location"
		fmt.Println("Please give proper command line arguments")
		return
	}
	startTimeArg := args[1]
	finishTimeArg := args[3]
	fileName := args[5]

	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println("cannot able to read the file", err)
		return
	}

	defer file.Close() //close after checking err

	queryStartTime, err := time.Parse("2006-01-02T15:04:05.0000Z", startTimeArg)
	if err != nil {
		fmt.Println("Could not able to parse the start time", startTimeArg)
		return
	}

	queryFinishTime, err := time.Parse("2006-01-02T15:04:05.0000Z", finishTimeArg)
	if err != nil {
		fmt.Println("Could not able to parse the finish time", finishTimeArg)
		return
	}

	filestat, err := file.Stat()
	if err != nil {
		fmt.Println("Could not able to get the file stat")
		return
	}

	fileSize := filestat.Size()
	offset := fileSize - 1
	lastLineSize := 0

	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "\n" {
			break
		}
		offset--
		lastLineSize += n
	}

	lastLine := make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset+1)

	if err != nil {
		fmt.Println("Could not able to read last line with offset", offset, "and lastline size", lastLineSize)
		return
	}

	logSlice := strings.SplitN(string(lastLine), ",", 2)
	logCreationTimeString := logSlice[0]

	lastLogCreationTime, err := time.Parse("2006-01-02T15:04:05.0000Z", logCreationTimeString)
	if err != nil {
		fmt.Println("can not able to parse time : ", err)
	}

	if lastLogCreationTime.After(queryStartTime) && lastLogCreationTime.Before(queryFinishTime) {
		Process(file, queryStartTime, queryFinishTime)
	}

	fmt.Println("\nTime taken - ", time.Since(s))
}

func Process(f *os.File, start time.Time, end time.Time) error {

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		wg.Add(1)
		go func() {
			ProcessChunk(buf, &linesPool, &stringPool, start, end)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, start time.Time, end time.Time) {

	var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {

		wg2.Add(1)
		go func(s int, e int) {
			defer wg2.Done() //to avaoid deadlocks
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				logSlice := strings.SplitN(text, ",", 2)
				logCreationTimeString := logSlice[0]

				logCreationTime, err := time.Parse("2006-01-02T15:04:05.0000Z", logCreationTimeString)
				if err != nil {
					fmt.Printf("\n Could not able to parse the time :%s for log : %v", logCreationTimeString, text)
					return
				}

				if logCreationTime.After(start) && logCreationTime.Before(end) {
					//fmt.Println(text)
				}
			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	wg2.Wait()
	logsSlice = nil
}
