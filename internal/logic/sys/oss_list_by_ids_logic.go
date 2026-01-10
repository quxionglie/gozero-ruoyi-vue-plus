// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strconv"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询OSS对象基于id串
func NewOssListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssListByIdsLogic {
	return &OssListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssListByIdsLogic) OssListByIds(req *types.OssListByIdsReq) (resp *types.OssListResp, err error) {
	// 1. 解析OSS ID列表
	ossIdStrs := strings.Split(req.OssIds, ",")
	var ossIds []int64
	for _, idStr := range ossIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue
		}
		ossIds = append(ossIds, id)
	}

	if len(ossIds) == 0 {
		return &types.OssListResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "OSS ID不能为空",
			},
			Total: 0,
			Rows:  []types.OssVo{},
		}, nil
	}

	// 2. 查询OSS对象列表
	ossList, err := l.svcCtx.SysOssModel.FindByIds(l.ctx, ossIds)
	if err != nil {
		l.Errorf("查询OSS对象列表失败: %v", err)
		return &types.OssListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询OSS对象列表失败",
			},
		}, err
	}

	// 3. 转换为响应格式
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
		Total: int64(len(rows)),
		Rows:  rows,
	}, nil
}

// convertToOssVo 转换OSS实体为响应格式
func (l *OssListByIdsLogic) convertToOssVo(oss *model.SysOss) types.OssVo {
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
