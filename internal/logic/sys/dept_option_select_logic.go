// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptOptionSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门选择框列表
func NewDeptOptionSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptOptionSelectLogic {
	return &DeptOptionSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptOptionSelectLogic) DeptOptionSelect(req *types.DeptOptionSelectReq) (resp *types.DeptListResp, err error) {
	var depts []*sys.SysDept

	// 1. 如果提供了 deptIds，按 ID 列表查询
	if req.DeptIds != "" {
		deptIdStrs := strings.Split(req.DeptIds, ",")
		var deptIds []int64
		for _, idStr := range deptIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			deptIds = append(deptIds, id)
		}
		if len(deptIds) > 0 {
			depts, err = l.svcCtx.SysDeptModel.FindByIds(l.ctx, deptIds)
			if err != nil {
				l.Errorf("查询部门列表失败: %v", err)
				return &types.DeptListResp{
					BaseResp: types.BaseResp{
						Code: 500,
						Msg:  "查询部门列表失败",
					},
				}, err
			}
		}
	}

	// 2. 转换为响应格式
	rows := make([]types.DeptVo, 0, len(depts))
	for _, dept := range depts {
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
