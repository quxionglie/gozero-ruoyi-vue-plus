// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询部门列表
func NewDeptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptListLogic {
	return &DeptListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptListLogic) DeptList(req *types.DeptListReq) (resp *types.DeptListResp, err error) {
	// 1. 构建查询条件
	deptQuery := &sys.DeptQuery{
		DeptId:       req.DeptId,
		ParentId:     req.ParentId,
		DeptName:     req.DeptName,
		DeptCategory: req.DeptCategory,
		Status:       req.Status,
		BelongDeptId: req.BelongDeptId,
	}

	// 2. 如果提供了 BelongDeptId，需要查询该部门及其所有子部门
	var depts []*sys.SysDept
	if req.BelongDeptId > 0 {
		// 查询该部门及其所有子部门的ID
		childDeptIds, err := l.svcCtx.SysDeptModel.SelectDeptAndChildById(l.ctx, req.BelongDeptId)
		if err == nil && len(childDeptIds) > 0 {
			// 使用 IN 查询
			allDepts, _ := l.svcCtx.SysDeptModel.FindAll(l.ctx, &sys.DeptQuery{
				DeptName:     req.DeptName,
				DeptCategory: req.DeptCategory,
				Status:       req.Status,
			})
			// 过滤出符合条件的部门
			for _, dept := range allDepts {
				for _, id := range childDeptIds {
					if dept.DeptId == id {
						depts = append(depts, dept)
						break
					}
				}
			}
		} else {
			// 如果查询失败，使用正常查询
			depts, err = l.svcCtx.SysDeptModel.FindAll(l.ctx, deptQuery)
		}
	} else {
		// 3. 正常查询数据
		depts, err = l.svcCtx.SysDeptModel.FindAll(l.ctx, deptQuery)
	}
	if err != nil {
		l.Errorf("查询部门列表失败: %v", err)
		return &types.DeptListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.DeptVo, 0, len(depts))
	for _, dept := range depts {
		deptVo := l.convertToDeptVo(dept)
		rows = append(rows, deptVo)
	}

	return &types.DeptListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rows,
	}, nil
}

// convertToDeptVo 转换部门实体为响应格式
func (l *DeptListLogic) convertToDeptVo(dept *sys.SysDept) types.DeptVo {
	deptVo := types.DeptVo{
		DeptId:       dept.DeptId,
		ParentId:     dept.ParentId,
		ParentName:   "",
		Ancestors:    dept.Ancestors,
		DeptName:     dept.DeptName,
		DeptCategory: "",
		OrderNum:     int32(dept.OrderNum),
		Leader:       0,
		LeaderName:   "",
		Phone:        "",
		Email:        "",
		Status:       dept.Status,
		CreateTime:   "",
		Children:     []types.DeptVo{},
	}

	if dept.DeptCategory.Valid {
		deptVo.DeptCategory = dept.DeptCategory.String
	}
	if dept.Leader.Valid {
		deptVo.Leader = dept.Leader.Int64
	}
	if dept.Phone.Valid {
		deptVo.Phone = dept.Phone.String
	}
	if dept.Email.Valid {
		deptVo.Email = dept.Email.String
	}
	if dept.CreateTime.Valid {
		deptVo.CreateTime = dept.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	// 查询父部门名称
	if dept.ParentId > 0 {
		parentDept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, dept.ParentId)
		if err == nil && parentDept != nil {
			deptVo.ParentName = parentDept.DeptName
		}
	}

	// 查询负责人名称（如果需要）
	if deptVo.Leader > 0 {
		// TODO: 从用户表查询负责人名称
	}

	return deptVo
}
