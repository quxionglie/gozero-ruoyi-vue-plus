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

type PostAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增岗位
func NewPostAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostAddLogic {
	return &PostAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostAddLogic) PostAdd(req *types.PostReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
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
	unique, err := l.svcCtx.SysPostModel.CheckPostNameUnique(l.ctx, req.PostName, req.DeptId, 0)
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
			Msg:  fmt.Sprintf("新增岗位'%s'失败，岗位名称已存在", req.PostName),
		}, nil
	}

	// 4. 校验岗位编码唯一性（全局唯一）
	unique, err = l.svcCtx.SysPostModel.CheckPostCodeUnique(l.ctx, req.PostCode, 0)
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
			Msg:  fmt.Sprintf("新增岗位'%s'失败，岗位编码已存在", req.PostName),
		}, nil
	}

	// 5. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 6. 生成主键ID（使用雪花算法）
	postId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成岗位ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成岗位ID失败",
		}, err
	}

	// 7. 构建岗位实体
	post := &model.SysPost{
		PostId:       postId,
		TenantId:     tenantId,
		DeptId:       req.DeptId,
		PostCode:     req.PostCode,
		PostCategory: sql.NullString{String: req.PostCategory, Valid: req.PostCategory != ""},
		PostName:     req.PostName,
		PostSort:     int64(req.PostSort),
		Status:       req.Status,
		Remark:       sql.NullString{String: req.Remark, Valid: req.Remark != ""},
		CreateDept:   sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		CreateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	if post.Status == "" {
		post.Status = "0"
	}
	if post.PostSort == 0 {
		post.PostSort = 0
	}

	// 8. 插入数据库
	_, err = l.svcCtx.SysPostModel.Insert(l.ctx, post)
	if err != nil {
		l.Errorf("新增岗位失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增岗位失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
