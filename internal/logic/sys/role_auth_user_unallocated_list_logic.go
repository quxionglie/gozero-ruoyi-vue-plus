// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAuthUserUnallocatedListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询未分配用户角色列表
func NewRoleAuthUserUnallocatedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAuthUserUnallocatedListLogic {
	return &RoleAuthUserUnallocatedListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleAuthUserUnallocatedListLogic) RoleAuthUserUnallocatedList(req *types.RoleAuthUserUnallocatedListReq) (resp *types.UserListResp, err error) {
	// 1. 构建查询条件
	userQuery := &model.UserQuery{
		RoleId:      req.RoleId,
		UserName:    req.UserName,
		Phonenumber: req.Phonenumber,
		Status:      req.Status,
	}

	// 2. 构建分页查询条件
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}
	if pageQuery.PageNum <= 0 {
		pageQuery.PageNum = 1
	}
	if pageQuery.PageSize <= 0 {
		pageQuery.PageSize = 10
	}

	// 3. 查询未分配用户列表
	userList, total, err := l.svcCtx.SysUserModel.FindUnallocatedPage(l.ctx, userQuery, pageQuery)
	if err != nil {
		l.Errorf("查询未分配用户列表失败: %v", err)
		return &types.UserListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询未分配用户列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.SysUserVo, 0, len(userList))
	for _, user := range userList {
		userVo := convertUserToVo(user)
		rows = append(rows, userVo)
	}

	return &types.UserListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Total: total,
		Rows:  rows,
	}, nil
}
