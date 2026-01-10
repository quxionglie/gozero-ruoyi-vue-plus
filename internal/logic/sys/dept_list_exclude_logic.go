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

type DeptListExcludeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询部门列表（排除节点）
func NewDeptListExcludeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptListExcludeLogic {
	return &DeptListExcludeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptListExcludeLogic) DeptListExclude(req *types.DeptListExcludeReq) (resp *types.DeptListResp, err error) {
	// 1. 查询所有部门
	deptQuery := &sys.DeptQuery{}
	depts, err := l.svcCtx.SysDeptModel.FindAll(l.ctx, deptQuery)
	if err != nil {
		l.Errorf("查询部门列表失败: %v", err)
		return &types.DeptListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门列表失败",
			},
		}, err
	}

	// 2. 排除指定部门及其所有子部门
	excludeDeptIdStr := strconv.FormatInt(req.DeptId, 10)
	var filteredDepts []*sys.SysDept
	for _, dept := range depts {
		// 排除指定部门本身
		if dept.DeptId == req.DeptId {
			continue
		}
		// 排除指定部门的所有子部门（ancestors 中包含该部门ID）
		ancestorsList := strings.Split(dept.Ancestors, ",")
		shouldExclude := false
		for _, ancestor := range ancestorsList {
			if strings.TrimSpace(ancestor) == excludeDeptIdStr {
				shouldExclude = true
				break
			}
		}
		if !shouldExclude {
			filteredDepts = append(filteredDepts, dept)
		}
	}

	// 3. 转换为响应格式
	rows := make([]types.DeptVo, 0, len(filteredDepts))
	for _, dept := range filteredDepts {
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
func (l *DeptListExcludeLogic) convertToDeptVo(dept *sys.SysDept) types.DeptVo {
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

	return deptVo
}
