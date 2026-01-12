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

type DictTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询字典类型列表
func NewDictTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeListLogic {
	return &DictTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeListLogic) DictTypeList(req *types.DictTypeListReq) (resp *types.TableDataInfoResp, err error) {
	// 构建查询条件
	dictTypeQuery := &model.DictTypeQuery{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
	}
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 使用 SQL 分页查询
	rows, total, err := l.svcCtx.SysDictTypeModel.FindPage(l.ctx, dictTypeQuery, pageQuery)
	if err != nil {
		l.Errorf("查询字典类型列表失败: %v", err)
		return &types.TableDataInfoResp{
			Total: 0,
			Rows:  []types.DictTypeVo{},
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询字典类型列表失败",
			},
		}, err
	}

	// 转换为 VO
	voList := make([]types.DictTypeVo, 0, len(rows))
	for _, row := range rows {
		vo := types.DictTypeVo{
			DictId:     row.DictId,
			DictName:   row.DictName,
			DictType:   row.DictType,
			Remark:     "",
			CreateTime: "",
		}
		if row.Remark.Valid {
			vo.Remark = row.Remark.String
		}
		if row.CreateTime.Valid {
			vo.CreateTime = row.CreateTime.Time.Format("2006-01-02 15:04:05")
		}
		voList = append(voList, vo)
	}

	return &types.TableDataInfoResp{
		Total: total,
		Rows:  voList,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}
