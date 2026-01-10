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

type PostDeptTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门树列表
func NewPostDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostDeptTreeLogic {
	return &PostDeptTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostDeptTreeLogic) PostDeptTree() (resp *types.DeptTreeResp, err error) {
	// 1. 获取租户ID（用于后续过滤，当前先不限制）
	_, err = util.GetTenantIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取租户ID失败: %v", err)
		return &types.DeptTreeResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取租户ID失败",
			},
		}, err
	}

	// 2. 构建查询条件
	deptQuery := &model.DeptQuery{
		// 查询所有部门
	}

	// 3. 查询所有部门
	deptList, err := l.svcCtx.SysDeptModel.FindAll(l.ctx, deptQuery)
	if err != nil {
		l.Errorf("查询部门列表失败: %v", err)
		return &types.DeptTreeResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门列表失败",
			},
		}, err
	}

	// 4. 过滤掉 deptName 为空的部门，并转换为 DeptTreeVo
	var filteredDeptList []*model.SysDept
	for _, dept := range deptList {
		if dept.DeptName != "" {
			filteredDeptList = append(filteredDeptList, dept)
		}
	}

	// 5. 找到所有根节点（parentId 不在任何节点的 deptId 中，或者 parentId <= 0）
	deptIdSet := make(map[int64]bool)
	for _, dept := range filteredDeptList {
		deptIdSet[dept.DeptId] = true
	}

	// 构建树形结构，将 DeptVo 转换为 DeptTreeVo
	var buildTree func(parentId int64) []types.DeptTreeVo
	buildTree = func(parentId int64) []types.DeptTreeVo {
		var children []types.DeptTreeVo
		for _, dept := range filteredDeptList {
			if dept.ParentId == parentId {
				disabled := dept.Status == "1" // status == "1" 表示禁用
				treeVo := types.DeptTreeVo{
					Id:       dept.DeptId,
					ParentId: dept.ParentId,
					Label:    dept.DeptName,
					Weight:   int32(dept.OrderNum),
					Disabled: disabled,
					Children: buildTree(dept.DeptId),
				}
				children = append(children, treeVo)
			}
		}
		return children
	}

	// 找到所有根节点（parentId 不存在于任何 deptId 中，或者 parentId <= 0）
	var rootNodes []types.DeptTreeVo
	for _, dept := range filteredDeptList {
		// 根节点：parentId <= 0 或者 parentId 不在任何 deptId 中
		if dept.ParentId <= 0 || !deptIdSet[dept.ParentId] {
			disabled := dept.Status == "1" // status == "1" 表示禁用
			treeVo := types.DeptTreeVo{
				Id:       dept.DeptId,
				ParentId: dept.ParentId,
				Label:    dept.DeptName,
				Weight:   int32(dept.OrderNum),
				Disabled: disabled,
				Children: buildTree(dept.DeptId),
			}
			rootNodes = append(rootNodes, treeVo)
		}
	}

	return &types.DeptTreeResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rootNodes,
	}, nil
}
