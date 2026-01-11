package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListByDeptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门下的所有用户信息
func NewUserListByDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListByDeptLogic {
	return &UserListByDeptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListByDeptLogic) UserListByDept(req *types.UserListByDeptReq) (resp *types.UserListByDeptResp, err error) {
	// 1. 查询部门下的所有用户
	users, err := l.svcCtx.SysUserModel.SelectUserListByDept(l.ctx, req.DeptId)
	if err != nil {
		l.Errorf("查询部门用户列表失败: %v", err)
		return &types.UserListByDeptResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门用户列表失败",
			},
		}, err
	}

	// 2. 转换为响应格式
	rows := make([]types.SysUserVo, 0, len(users))
	for _, user := range users {
		userVo := convertUserToVo(l.ctx, l.svcCtx, user)
		rows = append(rows, userVo)
	}

	return &types.UserListByDeptResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rows,
	}, nil
}
