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

type DeptGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询部门详细
func NewDeptGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptGetInfoLogic {
	return &DeptGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptGetInfoLogic) DeptGetInfo(req *types.DeptGetInfoReq) (resp *types.DeptResp, err error) {
	// 1. 查询部门信息
	dept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.DeptId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.DeptResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "部门不存在",
				},
			}, nil
		}
		l.Errorf("查询部门信息失败: %v", err)
		return &types.DeptResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门信息失败",
			},
		}, err
	}

	// 2. 转换为响应格式
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

	// 3. 查询父部门名称
	if dept.ParentId > 0 {
		parentDept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, dept.ParentId)
		if err == nil && parentDept != nil {
			deptVo.ParentName = parentDept.DeptName
		}
	}

	return &types.DeptResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: deptVo,
	}, nil
}
