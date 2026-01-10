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

type PostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询岗位列表
func NewPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PostListLogic {
	return &PostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PostListLogic) PostList(req *types.PostListReq) (resp *types.TableDataInfoResp, err error) {
	// 1. 设置默认分页参数
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 2. 构建查询条件
	postQuery := &sys.PostQuery{
		PostCode:     req.PostCode,
		PostCategory: req.PostCategory,
		PostName:     req.PostName,
		Status:       req.Status,
		DeptId:       req.DeptId,
		BelongDeptId: req.BelongDeptId,
	}
	pageQuery := &sys.PageQuery{
		PageNum:       pageNum,
		PageSize:      pageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 3. 查询数据
	posts, total, err := l.svcCtx.SysPostModel.FindPage(l.ctx, postQuery, pageQuery)
	if err != nil {
		l.Errorf("查询岗位列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询岗位列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
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
		// TODO: 查询部门名称（如果需要）
		rows = append(rows, postVo)
	}

	return &types.TableDataInfoResp{
		Total: total,
		Rows:  rows,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
