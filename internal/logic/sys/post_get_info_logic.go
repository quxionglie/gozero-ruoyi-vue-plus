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

type PostGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询岗位详细
func NewPostGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostGetInfoLogic {
	return &PostGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostGetInfoLogic) PostGetInfo(req *types.PostGetInfoReq) (resp *types.PostResp, err error) {
	// 1. 查询岗位信息
	post, err := l.svcCtx.SysPostModel.FindOne(l.ctx, req.PostId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.PostResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "岗位不存在",
				},
			}, nil
		}
		l.Errorf("查询岗位信息失败: %v", err)
		return &types.PostResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询岗位信息失败",
			},
		}, err
	}

	// 2. 转换为响应格式
	postVo := types.PostVo{
		PostId:       post.PostId,
		DeptId:       post.DeptId,
		PostCode:     post.PostCode,
		PostCategory: "",
		PostName:     post.PostName,
		PostSort:     int32(post.PostSort),
		Status:       post.Status,
		Remark:       "",
		CreateTime:   "",
		DeptName:     "",
	}
	if post.PostCategory.Valid {
		postVo.PostCategory = post.PostCategory.String
	}
	if post.Remark.Valid {
		postVo.Remark = post.Remark.String
	}
	if post.CreateTime.Valid {
		postVo.CreateTime = post.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	// TODO: 查询部门名称（如果需要）

	return &types.PostResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: postVo,
	}, nil
}
