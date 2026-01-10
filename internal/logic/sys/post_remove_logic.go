// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PostRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除岗位
func NewPostRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostRemoveLogic {
	return &PostRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostRemoveLogic) PostRemove(req *types.PostRemoveReq) (resp *types.BaseResp, err error) {
	if req.PostIds == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "岗位ID不能为空",
		}, nil
	}

	// 解析岗位ID列表（逗号分隔）
	postIdStrs := strings.Split(req.PostIds, ",")
	var postIds []int64
	for _, idStr := range postIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("岗位ID格式错误: %s", idStr),
			}, nil
		}
		postIds = append(postIds, id)
	}

	if len(postIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "岗位ID不能为空",
		}, nil
	}

	// 检查每个岗位是否已分配用户
	for _, postId := range postIds {
		count, err := l.svcCtx.SysPostModel.CountUserPostById(l.ctx, postId)
		if err != nil {
			l.Errorf("统计岗位使用数量失败: postId=%d, err=%v", postId, err)
			continue
		}
		if count > 0 {
			// 查询岗位名称
			post, err := l.svcCtx.SysPostModel.FindOne(l.ctx, postId)
			if err == nil {
				return &types.BaseResp{
					Code: 500,
					Msg:  fmt.Sprintf("%s已分配，不能删除", post.PostName),
				}, nil
			}
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("岗位ID %d 已分配，不能删除", postId),
			}, nil
		}
	}

	// 删除岗位
	for _, postId := range postIds {
		err = l.svcCtx.SysPostModel.Delete(l.ctx, postId)
		if err != nil {
			l.Errorf("删除岗位失败: postId=%d, err=%v", postId, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除岗位失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
