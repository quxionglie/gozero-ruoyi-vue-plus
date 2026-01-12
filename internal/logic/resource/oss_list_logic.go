// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询OSS对象存储列表
func NewOssListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssListLogic {
	return &OssListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssListLogic) OssList(req *types.OssListReq) (resp *types.OssListResp, err error) {
	// 1. 构建查询条件
	ossQuery := &model.OssQuery{
		FileName:     req.FileName,
		OriginalName: req.OriginalName,
		FileSuffix:   req.FileSuffix,
		Url:          req.Url,
		CreateBy:     req.CreateBy,
		Service:      req.Service,
	}

	// 2. 构建分页查询条件
	pageQuery := &model.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 3. 查询数据
	ossList, total, err := l.svcCtx.SysOssModel.FindPage(l.ctx, ossQuery, pageQuery)
	if err != nil {
		l.Errorf("查询OSS对象存储列表失败: %v", err)
		return &types.OssListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询OSS对象存储列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.OssVo, 0, len(ossList))
	for _, oss := range ossList {
		ossVo := l.convertToOssVo(oss)
		rows = append(rows, ossVo)
	}

	return &types.OssListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Total: total,
		Rows:  rows,
	}, nil
}

// convertToOssVo 转换OSS实体为响应格式
func (l *OssListLogic) convertToOssVo(oss *model.SysOss) types.OssVo {
	ossVo := types.OssVo{
		OssId:        oss.OssId,
		FileName:     oss.FileName,
		OriginalName: oss.OriginalName,
		FileSuffix:   oss.FileSuffix,
		Url:          oss.Url,
		Ext1:         "",
		CreateTime:   "",
		CreateBy:     0,
		CreateByName: "",
		Service:      oss.Service,
	}

	if oss.Ext1.Valid {
		ossVo.Ext1 = oss.Ext1.String
	}
	if oss.CreateTime.Valid {
		ossVo.CreateTime = oss.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	if oss.CreateBy.Valid {
		ossVo.CreateBy = oss.CreateBy.Int64
		// TODO: 查询用户名
	}

	return ossVo
}
