// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据角色编号获取详细信息
func NewRoleGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleGetInfoLogic {
	return &RoleGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleGetInfoLogic) RoleGetInfo(req *types.RoleGetInfoReq) (resp *types.RoleResp, err error) {
	// 1. 查询角色信息
	role, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, req.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.RoleResp{
				BaseResp: types.BaseResp{
					Code: 404,
					Msg:  "角色不存在",
				},
			}, nil
		}
		l.Errorf("查询角色信息失败: %v", err)
		return &types.RoleResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询角色信息失败",
			},
		}, err
	}

	// 2. 查询角色菜单ID列表
	menuIds, err := l.svcCtx.SysMenuModel.SelectMenuListByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("查询角色菜单列表失败: %v", err)
		menuIds = []int64{}
	}

	// 3. 查询角色部门ID列表（数据权限）
	deptIds, err := l.svcCtx.SysRoleDeptModel.SelectDeptIdsByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("查询角色部门列表失败: %v", err)
		deptIds = []int64{}
	}

	// 4. 转换为响应格式
	roleVo := convertRoleToVo(role)
	roleVo.MenuIds = util.Int64SliceToStringSlice(menuIds)
	roleVo.DeptIds = util.Int64SliceToStringSlice(deptIds)

	return &types.RoleResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: roleVo,
	}, nil
}
