package main

//redis的分布式锁，用的格式 setNX命令
//set命令：set key value ex seconds nx    set  lock  11   ex  5  nx
//ex 表示过期时间，精确到秒 （对应另一个参数px过期时间精确到毫秒)
//nx 表示if not exists，只有键不存在才能设置成功（对应另一个参数xx只有键存在才能设置成功）
//https://github.com/skyhackvip/lock/blob/master/redislock.go
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

//创建redis的客户端
var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

//数据加锁，初始方案
func lock(myfunc func()) {

	var lockKey = "mylocker"

	//上锁
	for {
		locksuccess, err := client.SetNX(context.Background(), lockKey, 1, time.Second*5).Result() //value:1  过期时间5秒
		if err != nil || !locksuccess {
			fmt.Println("get lock failed")
			time.Sleep(time.Microsecond * 10)
			//这里最好停留一会，防止自旋太频繁
			continue
		} else {
			fmt.Println("lock success")
			break
		}
	}

	//TODO:
	//做自己的逻辑代码
	myfunc()

	//解锁
	_, err := client.Del(context.Background(), lockKey).Result()
	if err != nil {
		fmt.Println("unlock failed")
		return
	} else {
		fmt.Println("unlock sucessed")
	}
}

var counter int64

func myfunc() {
	counter++
	fmt.Println(counter)
}

var wg sync.WaitGroup

func main() {
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lock(myfunc)
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

//结果
/*
get lock failed
get lock failed
lock success
get lock failed
get lock failed
unlock sucessed
1

多个线程抢锁，只有一个执行完毕，其他全部失败
*/

//优化1：使用uuid标识锁
//1的代码有问题：
//上一份代码中，如果 G1 执行任务的时间超过了过期时间，那么就会自动解锁
//而此时G2 又或得到了锁，当G1 正真运行到解锁时，解锁的是G2上的锁
//当G2解锁时，就会报错，尝试解锁一个未上锁的锁
//方案优化，为每个goroutine设置一个uuid，解锁的时候根据uuid进行解锁，如果不匹配，说明当前goroutine的锁过期了
func lockV2(myfunc func()) {
	id, _ := uuid.NewV4()
	var lockKey = "mylocker"
	locksuccess, err := client.SetNX(context.Background(), lockKey, id, time.Second*2).Result()
	if err != nil || !locksuccess {
		fmt.Println("get lock failed")
	} else {
		fmt.Println("lock success")
	}

	myfunc()

	//解锁，先判断当前锁是不是属于当前线程的了
	v, _ := client.Get(context.Background(), lockKey).Result()
	if v == id.String() {
		//锁还在，那么就解锁
		_, err = client.Del(context.Background(), lockKey).Result()
	} else {
		fmt.Println("锁已经超时了")
	}

}

//优化2：使用lua脚本
//v2还是存在问题：获取v 、 判断 v==uuid 、 解锁这三步不是原子操作，还是会出错
//解决办法：利用lua脚本
//lua是嵌入式语言，redis本身支持。使用golang操作redis运行lua命令，保障问题解决。
//lua脚本中KEYS[1]代表lock的key，ARGV[1]代表lock的value，也就是生成的uuid。通过执行lua来保障这里删除锁的操作是原子的。
//完整代码参见：https://github.com/skyhackvip/lock/blob/master/redislualock.go
func lockV3(myfunc func()) {
	id, _ := uuid.NewV4()
	var lockKey = "mylocker"
	locksuccess, err := client.SetNX(context.Background(), lockKey, id, time.Second*2).Result()
	if err != nil || !locksuccess {
		fmt.Println("get lock failed")
	} else {
		fmt.Println("lock success")
	}

	myfunc()
	//使用lua脚本进行解锁
	//解锁失败返回0
	var luaStript = redis.NewScript(`
		if redis.call("get",KEYS[1]) == ARGV[1]
			then 
				return redis.call("del",KEYS[1])
			else
				return 0
		end
	`)

	rs, err := luaStript.Run(context.Background(), client, []string{lockKey}, id).Result()
	if rs == 0 {
		fmt.Println("unlock failed")
	} else {
		fmt.Println("unlock successed")
	}
}
