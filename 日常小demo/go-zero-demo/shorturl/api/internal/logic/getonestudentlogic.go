package logic

import (
	"context"

	"shorturl/api/internal/svc"
	"shorturl/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetOneStudentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOneStudentLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOneStudentLogic {
	return GetOneStudentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOneStudentLogic) GetOneStudent(req types.GetonestudentReq) (*types.GetonestudentResp, error) {
	// todo: add your logic here and delete this line

	return &types.GetonestudentResp{
		Name: req.Name,
		Age:  18,
	}, nil
}
