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

type RoleDeptTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对应角色部门树列表
func NewRoleDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDeptTreeLogic {
	return &RoleDeptTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDeptTreeLogic) RoleDeptTree(req *types.RoleDeptTreeReq) (resp *types.DeptTreeSelectResp, err error) {
	// 1. 查询角色关联的部门ID列表
	deptIds, err := l.svcCtx.SysRoleDeptModel.SelectDeptIdsByRoleId(l.ctx, req.RoleId)
	if err != nil {
		l.Errorf("查询角色部门列表失败: %v", err)
		return &types.DeptTreeSelectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询角色部门列表失败",
			},
		}, err
	}

	// 2. 查询所有部门（用于构建部门树）
	depts, err := l.svcCtx.SysDeptModel.FindAll(l.ctx, &model.DeptQuery{})
	if err != nil {
		l.Errorf("查询部门列表失败: %v", err)
		return &types.DeptTreeSelectResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门列表失败",
			},
		}, err
	}

	// 3. 构建部门树
	deptTreeList := l.buildDeptTreeSelect(depts)

	return &types.DeptTreeSelectResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.DeptTreeSelectVo{
			CheckedKeys: util.Int64SliceToStringSlice(deptIds),
			Depts:       deptTreeList,
		},
	}, nil
}

// buildDeptTreeSelect 构建部门树形结构
func (l *RoleDeptTreeLogic) buildDeptTreeSelect(depts []*model.SysDept) []types.DeptTreeVo {
	if len(depts) == 0 {
		return []types.DeptTreeVo{}
	}

	// 构建树形结构
	var buildTree func(parentId int64) []types.DeptTreeVo
	buildTree = func(parentId int64) []types.DeptTreeVo {
		var children []types.DeptTreeVo
		for _, dept := range depts {
			if dept.ParentId == parentId {
				treeNode := types.DeptTreeVo{
					Id:       dept.DeptId,
					Label:    dept.DeptName,
					ParentId: dept.ParentId,
					Children: buildTree(dept.DeptId),
				}
				children = append(children, treeNode)
			}
		}
		return children
	}

	return buildTree(0)
}
