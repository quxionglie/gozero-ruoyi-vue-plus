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

type PostOptionSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取岗位选择框列表
func NewPostOptionSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostOptionSelectLogic {
	return &PostOptionSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostOptionSelectLogic) PostOptionSelect(req *types.PostOptionSelectReq) (resp *types.PostOptionSelectResp, err error) {
	var posts []*sys.SysPost

	// 1. 如果提供了 deptId，按部门查询
	if req.DeptId > 0 {
		postQuery := &sys.PostQuery{
			DeptId: req.DeptId,
		}
		posts, err = l.svcCtx.SysPostModel.FindAll(l.ctx, postQuery)
		if err != nil {
			l.Errorf("查询岗位列表失败: %v", err)
			return &types.PostOptionSelectResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "查询岗位列表失败",
				},
			}, err
		}
	} else if req.PostIds != "" {
		// 2. 如果提供了 postIds，按 ID 列表查询
		postIdStrs := strings.Split(req.PostIds, ",")
		var postIds []int64
		for _, idStr := range postIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				continue
			}
			postIds = append(postIds, id)
		}
		if len(postIds) > 0 {
			posts, err = l.svcCtx.SysPostModel.FindByIds(l.ctx, postIds)
			if err != nil {
				l.Errorf("查询岗位列表失败: %v", err)
				return &types.PostOptionSelectResp{
					BaseResp: types.BaseResp{
						Code: 500,
						Msg:  "查询岗位列表失败",
					},
				}, err
			}
		}
	}

	// 3. 转换为响应格式
	rows := make([]types.PostVo, 0, len(posts))
	for _, post := range posts {
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
		rows = append(rows, postVo)
	}

	return &types.PostOptionSelectResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: rows,
	}, nil
}
