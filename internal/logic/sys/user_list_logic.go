package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户列表
func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListReq) (resp *types.UserListResp, err error) {
	// 1. 获取租户ID
	_, err = util.GetTenantIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取租户ID失败: %v", err)
		return &types.UserListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取租户ID失败",
			},
		}, err
	}

	// 2. 构建查询条件
	userQuery := &model.UserQuery{
		UserId:         req.UserId,
		UserIds:        req.UserIds,
		UserName:       req.UserName,
		NickName:       req.NickName,
		Status:         req.Status,
		Phonenumber:    req.Phonenumber,
		DeptId:         req.DeptId,
		BeginTime:      req.BeginTime,
		EndTime:        req.EndTime,
		ExcludeUserIds: req.ExcludeUserIds,
	}

	// 3. 构建分页查询条件
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

	// 4. 查询用户列表
	userList, total, err := l.svcCtx.SysUserModel.FindPage(l.ctx, userQuery, pageQuery)
	if err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return &types.UserListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户列表失败",
			},
		}, err
	}

	// 5. 转换为响应格式
	rows := make([]types.SysUserVo, 0, len(userList))
	for _, user := range userList {
		userVo := convertUserToVo(l.ctx, l.svcCtx, user)
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
