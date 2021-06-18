package mycache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/零碎知识/mycache/consistenthash"
	pb "github.com/零碎知识/mycache/mycachepb"
	"google.golang.org/protobuf/proto"
)

//分布式缓存需要实现节点间通信，建立基于 HTTP 的通信机制是比较常见和简单的做法
//如果一个节点启动了 HTTP 服务，那么这个节点就可以被其他节点访问。
//单机节点构造HTTP服务，不与其他部分耦合，提供被其它节点访问的能力

//首先我们创建一个结构体 HTTPPool，作为承载节点间 HTTP 通信的核心数据结构（包括服务端和客户端，今天只实现服务端）。

const defaultBasePath = "/_mycache/"
const defaultReplices = 50

type HTTPPool struct {
	self     string //用来记录自己的 ip地址和端口号
	basePath string //处理 /basepath/的请求，因为考虑到本机上有其他服务，如大部分网站API接口是以 /api 作为前缀
	//提供选择节点的能力
	mu          sync.Mutex
	peers       *consistenthash.MyMap  //类型是一致哈希Map，用来根据具体的key选择节点
	httpGetters map[string]*httpGetter //映射远程节点与对应的httpGetter，每个远程节点对应httpGetter，因为 httpGetter 与远程节点的地址 baseURL 有关。
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:        self,
		basePath:    defaultBasePath,
		httpGetters: make(map[string]*httpGetter),
	}
}

func (this *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.self, fmt.Sprintf(format, v...))
}

//然后实现ServeHTTP方法   ServeHTTP handle all http requests
func (this *HTTPPool) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//先判断访问前缀是不是符合要求
	if !strings.HasPrefix(r.URL.Path, this.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	this.Log("%s %s", r.Method, r.URL.Path)

	//定义，访问缓存必须是这样的格式 // <basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(this.basePath):], "/", 2) // 分成最多2个子串; 考虑到万一key里面也有 / 这个符号
	if len(parts) != 2 {
		http.Error(rw, "bad request", http.StatusBadRequest)
		return
	}

	//接下来就获取到了 group 和 key 了
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(rw, "no such group:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		//没获取到，算是内部错误
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the value to the response body as a proto message.
	// 将value写进response body，作为一个proto的消息
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//到这才算获取到了数据
	//	w.Header().Set("Content-Type", "application/octet-stream") 这是以二进制文件的形式发送，不太对
	//rw.Write(view.ByteSlice()) //直接传到前端
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Write(body)

}

//Set更新httppool中的peers
func (this *HTTPPool) Set(peers ...string) {

	this.mu.Lock()
	defer this.mu.Unlock()
	this.peers = consistenthash.NewMap(defaultReplices, nil)
	this.peers.Add(peers...)     //增加节点
	for _, peer := range peers { //设置每个节点都有一个对应的httpGetter结构体，可以调用Get()函数，得到想要的数据
		this.httpGetters[peer] = &httpGetter{baseURL: peer + this.basePath}
	}

}

//PickPeer 根据key值挑选peer节点
func (this *HTTPPool) PickPeer(key string) (PeerGetter, bool) {

	this.mu.Lock()
	defer this.mu.Unlock()
	if peer := this.peers.Get(key); peer != "" || peer != this.self {
		this.Log("Pick peer %s", peer)
		return this.httpGetters[peer], true
	}
	return nil, false
}

//上面的 HTTPPool实现的是服务端的功能，接下来实现客户端的功能
//因为HTTPPool要从远程节点获取数据，所以得具备客户端的功能
type httpGetter struct {
	baseURL string
	//baseURL 表示将要访问的远程节点的地址，例如 http://example.com/_geecache/
}

func (this *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v/%v/%v",
		this.baseURL,
		url.QueryEscape(in.GetGroup()), //作用，如果group=a/ a，正常放进url肯定会出错，经过处理之后变成  a%2F+a,将 / 转义成 %2F,空格转义成 +
		url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	//从resp中读取数据
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("response has not ok status")
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ PeerGetter = (*httpGetter)(nil)

//测试代码
/*
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/零碎知识/mycache"
)

var db = map[string]string{
	"sjw": "123",
	"tao": "234",
	"jse": "345",
}

func main() {
	//先得创建自己group的回调函数
	mycache.NewGroup("score", 1024*10, mycache.GetterFunc(
		func(key string) ([]byte, error) {

			log.Println("[Slow DB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		},
	))

	addr := "localhost:6789"
	peers := mycache.NewHTTPPool(addr)
	fmt.Println("mycache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

*/
