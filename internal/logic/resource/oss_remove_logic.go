// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"strconv"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除OSS对象存储
func NewOssRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssRemoveLogic {
	return &OssRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssRemoveLogic) OssRemove(req *types.OssRemoveReq) (resp *types.BaseResp, err error) {
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
		return &types.BaseResp{
			Code: 400,
			Msg:  "OSS ID不能为空",
		}, nil
	}

	// 2. 查询OSS对象信息（用于删除OSS服务中的文件）
	ossList, err := l.svcCtx.SysOssModel.FindByIds(l.ctx, ossIds)
	if err != nil {
		l.Errorf("查询OSS对象信息失败: %v", err)
		// 即使查询失败，也继续删除数据库记录
	}

	// 3. 删除OSS服务中的文件
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)
	for _, ossObj := range ossList {
		if ossObj.Service != "" {
			ossClient, err := l.svcCtx.OssManager.GetClientByConfigKey(l.ctx, ossObj.Service, tenantId)
			if err != nil {
				l.Infof("获取OSS客户端失败（跳过删除OSS文件）: service=%s, error=%v", ossObj.Service, err)
				continue
			}

			// 删除OSS服务中的文件
			err = ossClient.Delete(ossObj.FileName)
			if err != nil {
				l.Infof("删除OSS服务中的文件失败（继续删除数据库记录）: fileName=%s, error=%v", ossObj.FileName, err)
				// 即使OSS文件删除失败，也继续删除数据库记录
			}
		}
	}

	// 4. 删除数据库记录
	for _, ossId := range ossIds {
		err = l.svcCtx.SysOssModel.Delete(l.ctx, ossId)
		if err != nil {
			if err == model.ErrNotFound {
				// 记录不存在，继续删除其他
				continue
			}
			l.Errorf("删除OSS对象失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除OSS对象失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
