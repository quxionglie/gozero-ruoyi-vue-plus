package sys

import (
	"context"
	"fmt"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOptionSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户ID串批量获取用户基础信息
func NewUserOptionSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOptionSelectLogic {
	return &UserOptionSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserOptionSelectLogic) UserOptionSelect(req *types.UserOptionSelectReq) (resp *types.UserOptionSelectResp, err error) {
	// 1. 解析用户ID串
	var userIds []int64
	if req.UserIds != "" {
		userIdsStr := strings.Split(req.UserIds, ",")
		for _, idStr := range userIdsStr {
			if idStr == "" {
				continue
			}
			var id int64
			_, err := fmt.Sscanf(strings.TrimSpace(idStr), "%d", &id)
			if err != nil {
				continue
			}
			userIds = append(userIds, id)
		}
	}

	// 2. 查询用户列表
	var users []*model.SysUser
	if len(userIds) > 0 {
		// 根据用户ID列表查询
		users, err = l.svcCtx.SysUserModel.FindByIds(l.ctx, userIds, req.DeptId)
	} else if req.DeptId > 0 {
		// 根据部门ID查询
		users, err = l.svcCtx.SysUserModel.SelectUserListByDept(l.ctx, req.DeptId)
	} else {
		// 查询所有正常状态的用户
		users, err = l.svcCtx.SysUserModel.FindByIds(l.ctx, []int64{}, 0)
	}

	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return &types.UserOptionSelectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户列表失败",
			},
		}, err
	}

	// 3. 转换为响应格式（只返回基础信息：userId, userName, nickName）
	rows := make([]types.SysUserVo, 0, len(users))
	for _, user := range users {
		// 只返回正常状态的用户
		if user.Status == "0" && user.DelFlag == "0" {
			userVo := convertUserToVo(user)
			rows = append(rows, userVo)
		}
	}

	return &types.UserOptionSelectResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rows,
	}, nil
}
