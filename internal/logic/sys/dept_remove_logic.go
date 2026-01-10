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

type DeptRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除部门
func NewDeptRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeptRemoveLogic {
	return &DeptRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeptRemoveLogic) DeptRemove(req *types.DeptRemoveReq) (resp *types.BaseResp, err error) {
	// 1. 默认部门（deptId=100）不允许删除
	if req.DeptId == 100 {
		return &types.BaseResp{
			Code: 500,
			Msg:  "默认部门,不允许删除",
		}, nil
	}

	// 2. 检查是否存在子部门
	hasChild, err := l.svcCtx.SysDeptModel.HasChildByDeptId(l.ctx, req.DeptId)
	if err != nil {
		l.Errorf("检查是否存在子部门失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查是否存在子部门失败",
		}, err
	}
	if hasChild {
		return &types.BaseResp{
			Code: 500,
			Msg:  "存在下级部门,不允许删除",
		}, nil
	}

	// 3. 检查部门是否存在用户
	existUser, err := l.svcCtx.SysDeptModel.CheckDeptExistUser(l.ctx, req.DeptId)
	if err != nil {
		l.Errorf("检查部门是否存在用户失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "检查部门是否存在用户失败",
		}, err
	}
	if existUser {
		return &types.BaseResp{
			Code: 500,
			Msg:  "部门存在用户,不允许删除",
		}, nil
	}

	// 4. 检查部门是否存在岗位
	postCount, err := l.svcCtx.SysPostModel.CountPostByDeptId(l.ctx, req.DeptId)
	if err != nil {
		l.Errorf("统计部门岗位数量失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "统计部门岗位数量失败",
		}, err
	}
	if postCount > 0 {
		return &types.BaseResp{
			Code: 500,
			Msg:  "部门存在岗位,不允许删除",
		}, nil
	}

	// 5. 删除部门
	err = l.svcCtx.SysDeptModel.Delete(l.ctx, req.DeptId)
	if err != nil {
		if err == sys.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "部门不存在",
			}, nil
		}
		l.Errorf("删除部门失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "删除部门失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
