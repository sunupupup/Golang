//抽象PeerPicker
package mycache

import pb "github.com/零碎知识/mycache/mycachepb"

//根据key值，求得存放这个数据的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

//从相应的 group中根据key值，获取value缓存
//type PeerGetter interface{
//	Get(group string,key string)([]byte,error)
//}
//修改 peers.go 中的 PeerGetter 接口，参数使用 geecachepb.pb.go 中的数据类型。
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
