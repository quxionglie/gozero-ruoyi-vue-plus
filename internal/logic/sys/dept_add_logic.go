// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeptAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增部门
func NewDeptAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptAddLogic {
	return &DeptAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptAddLogic) DeptAdd(req *types.DeptReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.DeptName == "" {
		return &types.BaseResp{Code: 400, Msg: "部门名称不能为空"}, nil
	}

	// 2. 参数长度校验
	if err := util.ValidateStringLength(req.DeptName, "部门名称", 30); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if req.Phone != "" {
		if err := util.ValidateStringLength(req.Phone, "联系电话", 11); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}
	if req.Email != "" {
		if err := util.ValidateStringLength(req.Email, "邮箱", 50); err != nil {
			return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
		}
	}

	// 3. 默认父部门ID为0（根部门）
	if req.ParentId == 0 {
		req.ParentId = 0
	}

	// 4. 校验部门名称唯一性（同父部门下唯一）
	unique, err := l.svcCtx.SysDeptModel.CheckDeptNameUnique(l.ctx, req.DeptName, req.ParentId, 0)
	if err != nil {
		l.Errorf("校验部门名称唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验部门名称唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("新增部门'%s'失败，部门名称已存在", req.DeptName),
		}, nil
	}

	// 5. 如果父节点不为正常状态,则不允许新增子节点
	if req.ParentId > 0 {
		parentDept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
		if err != nil {
			if err == model.ErrNotFound {
				return &types.BaseResp{
					Code: 500,
					Msg:  "父部门不存在",
				}, nil
			}
			l.Errorf("查询父部门失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "查询父部门失败",
			}, err
		}
		if parentDept.Status != "0" {
			return &types.BaseResp{
				Code: 500,
				Msg:  "部门停用，不允许新增",
			}, nil
		}
	}

	// 6. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 7. 生成主键ID（使用雪花算法）
	newDeptId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成部门ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成部门ID失败",
		}, err
	}

	// 8. 构建 ancestors（祖级列表）
	var ancestors string
	if req.ParentId == 0 {
		ancestors = "0"
	} else {
		parentDept, _ := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.ParentId)
		if parentDept != nil {
			ancestors = parentDept.Ancestors + "," + fmt.Sprintf("%d", req.ParentId)
		} else {
			ancestors = "0," + fmt.Sprintf("%d", req.ParentId)
		}
	}

	// 9. 构建部门实体
	dept := &model.SysDept{
		DeptId:       newDeptId,
		TenantId:     tenantId,
		ParentId:     req.ParentId,
		Ancestors:    ancestors,
		DeptName:     req.DeptName,
		DeptCategory: sql.NullString{String: req.DeptCategory, Valid: req.DeptCategory != ""},
		OrderNum:     int64(req.OrderNum),
		Leader:       sql.NullInt64{Int64: req.Leader, Valid: req.Leader > 0},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Email:        sql.NullString{String: req.Email, Valid: req.Email != ""},
		Status:       req.Status,
		DelFlag:      "0",
		CreateDept:   sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		CreateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	if dept.Status == "" {
		dept.Status = "0"
	}
	if dept.OrderNum == 0 {
		dept.OrderNum = 0
	}

	// 10. 插入数据库
	_, err = l.svcCtx.SysDeptModel.Insert(l.ctx, dept)
	if err != nil {
		l.Errorf("新增部门失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增部门失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
