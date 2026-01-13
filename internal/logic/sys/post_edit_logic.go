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

type PostEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改岗位
func NewPostEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostEditLogic {
	return &PostEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostEditLogic) PostEdit(req *types.PostReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.PostId <= 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "岗位ID不能为空",
		}, nil
	}
	if req.PostCode == "" {
		return &types.BaseResp{Code: 400, Msg: "岗位编码不能为空"}, nil
	}
	if req.PostName == "" {
		return &types.BaseResp{Code: 400, Msg: "岗位名称不能为空"}, nil
	}

	// 2. 参数长度校验
	if err := util.ValidateStringLength(req.PostCode, "岗位编码", 64); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}
	if err := util.ValidateStringLength(req.PostName, "岗位名称", 50); err != nil {
		return &types.BaseResp{Code: 400, Msg: err.Error()}, nil
	}

	// 3. 校验岗位名称唯一性（同部门内唯一）
	unique, err := l.svcCtx.SysPostModel.CheckPostNameUnique(l.ctx, req.PostName, req.DeptId, req.PostId)
	if err != nil {
		l.Errorf("校验岗位名称唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验岗位名称唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改岗位'%s'失败，岗位名称已存在", req.PostName),
		}, nil
	}

	// 4. 校验岗位编码唯一性（全局唯一）
	unique, err = l.svcCtx.SysPostModel.CheckPostCodeUnique(l.ctx, req.PostCode, req.PostId)
	if err != nil {
		l.Errorf("校验岗位编码唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验岗位编码唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改岗位'%s'失败，岗位编码已存在", req.PostName),
		}, nil
	}

	// 5. 检查岗位是否存在
	post, err := l.svcCtx.SysPostModel.FindOne(l.ctx, req.PostId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "岗位不存在",
			}, nil
		}
		l.Errorf("查询岗位失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询岗位失败",
		}, err
	}

	// 6. 如果禁用且已分配用户，不能禁用
	if req.Status == "1" && post.Status == "0" {
		count, err := l.svcCtx.SysPostModel.CountUserPostById(l.ctx, req.PostId)
		if err != nil {
			l.Errorf("统计岗位使用数量失败: %v", err)
		} else if count > 0 {
			return &types.BaseResp{
				Code: 500,
				Msg:  "该岗位下存在已分配用户，不能禁用",
			}, nil
		}
	}

	// 7. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 8. 更新岗位信息（只设置表单输入的字段）
	updatePost := &model.SysPost{
		PostId:     req.PostId,
		PostCode:   req.PostCode,
		PostName:   req.PostName,
		UpdateBy:   sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime: sql.NullTime{Time: time.Now(), Valid: true},
	}
	if req.DeptId > 0 {
		updatePost.DeptId = req.DeptId
	}
	if req.PostCategory != "" {
		updatePost.PostCategory = sql.NullString{String: req.PostCategory, Valid: true}
	}
	if req.PostSort > 0 {
		updatePost.PostSort = int64(req.PostSort)
	}
	if req.Status != "" {
		updatePost.Status = req.Status
	}
	if req.Remark != "" {
		updatePost.Remark = sql.NullString{String: req.Remark, Valid: true}
	}

	// 9. 更新数据库
	err = l.svcCtx.SysPostModel.UpdateById(l.ctx, updatePost)
	if err != nil {
		l.Errorf("修改岗位失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改岗位失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
