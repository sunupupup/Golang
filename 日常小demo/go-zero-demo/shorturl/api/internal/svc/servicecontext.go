package svc

import (
	"shorturl/api/internal/config"
	"shorturl/rpc/transform/transformer"

	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	//transformer.Transformer 是 ransformer rpc 服务对外暴露的接口
	Transformer transformer.Transformer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		//zrpc.MustNewClient(c.Transform) 创建了一个 grpc 客户端，用这个客户端调用远程函数
		Transformer: transformer.NewTransformer(zrpc.MustNewClient(c.Transform)), // 手动代码
	}
}
