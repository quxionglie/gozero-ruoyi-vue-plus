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

type DeptEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改部门
func NewDeptEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptEditLogic {
	return &DeptEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptEditLogic) DeptEdit(req *types.DeptReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.DeptId <= 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "部门ID不能为空",
		}, nil
	}
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

	// 3. 检查部门是否存在
	dept, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.DeptId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "部门不存在",
			}, nil
		}
		l.Errorf("查询部门失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询部门失败",
		}, err
	}

	// 4. 默认父部门ID为0（根部门）
	if req.ParentId == 0 {
		req.ParentId = 0
	}

	// 5. 上级部门不能是自己
	if req.ParentId == req.DeptId {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改部门'%s'失败，上级部门不能是自己", req.DeptName),
		}, nil
	}

	// 6. 校验部门名称唯一性（同父部门下唯一）
	unique, err := l.svcCtx.SysDeptModel.CheckDeptNameUnique(l.ctx, req.DeptName, req.ParentId, req.DeptId)
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
			Msg:  fmt.Sprintf("修改部门'%s'失败，部门名称已存在", req.DeptName),
		}, nil
	}

	// 7. 如果禁用，需要检查是否有子部门，是否有用户
	if req.Status == "1" && dept.Status == "0" {
		// 检查是否有正常状态的子部门
		count, err := l.svcCtx.SysDeptModel.CountNormalChildrenDeptById(l.ctx, req.DeptId)
		if err != nil {
			l.Errorf("统计子部门数量失败: %v", err)
		} else if count > 0 {
			return &types.BaseResp{
				Code: 500,
				Msg:  "该部门包含未停用的子部门!",
			}, nil
		}

		// 检查是否有用户
		existUser, err := l.svcCtx.SysDeptModel.CheckDeptExistUser(l.ctx, req.DeptId)
		if err != nil {
			l.Errorf("检查部门是否存在用户失败: %v", err)
		} else if existUser {
			return &types.BaseResp{
				Code: 500,
				Msg:  "该部门下存在已分配用户，不能禁用!",
			}, nil
		}
	}

	// 8. 如果父部门改变，需要更新 ancestors
	var ancestors string
	if req.ParentId != dept.ParentId {
		if req.ParentId == 0 {
			ancestors = "0"
		} else {
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
			ancestors = parentDept.Ancestors + "," + fmt.Sprintf("%d", req.ParentId)
		}
	} else {
		ancestors = dept.Ancestors
	}

	// 9. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 10. 更新部门信息（只设置表单输入的字段）
	updateDept := &model.SysDept{
		DeptId:     req.DeptId,
		DeptName:   req.DeptName,
		Ancestors:  ancestors,
		UpdateBy:   sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime: sql.NullTime{Time: time.Now(), Valid: true},
	}
	if req.ParentId > 0 {
		updateDept.ParentId = req.ParentId
	}
	if req.DeptCategory != "" {
		updateDept.DeptCategory = sql.NullString{String: req.DeptCategory, Valid: true}
	}
	if req.OrderNum > 0 {
		updateDept.OrderNum = int64(req.OrderNum)
	}
	if req.Leader > 0 {
		updateDept.Leader = sql.NullInt64{Int64: req.Leader, Valid: true}
	}
	if req.Phone != "" {
		updateDept.Phone = sql.NullString{String: req.Phone, Valid: true}
	}
	if req.Email != "" {
		updateDept.Email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.Status != "" {
		updateDept.Status = req.Status
	}

	// 11. 更新数据库
	err = l.svcCtx.SysDeptModel.UpdateById(l.ctx, updateDept)
	if err != nil {
		l.Errorf("修改部门失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改部门失败",
		}, err
	}

	// 12. TODO: 如果父部门改变，需要更新所有子部门的 ancestors
	// TODO: 如果部门状态为启用，需要启用所有上级部门

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
