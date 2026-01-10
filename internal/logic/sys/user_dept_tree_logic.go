package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeptTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取部门树列表
func NewUserDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeptTreeLogic {
	return &UserDeptTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserDeptTreeLogic) UserDeptTree() (resp *types.DeptListResp, err error) {
	// 1. 获取租户ID（用于后续过滤，当前先不限制）
	_, err = util.GetTenantIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取租户ID失败: %v", err)
		return &types.DeptListResp{
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
		return &types.DeptListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询部门列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.DeptVo, 0, len(deptList))
	for _, dept := range deptList {
		deptVo := types.DeptVo{
			DeptId:     dept.DeptId,
			ParentId:   dept.ParentId,
			DeptName:   dept.DeptName,
			OrderNum:   int32(dept.OrderNum),
			Status:     dept.Status,
			ParentName: "",
			Ancestors:  dept.Ancestors,
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
		if dept.DeptCategory.Valid {
			deptVo.DeptCategory = dept.DeptCategory.String
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
