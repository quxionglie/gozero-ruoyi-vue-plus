// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssConfigRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除对象存储配置
func NewOssConfigRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigRemoveLogic {
	return &OssConfigRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigRemoveLogic) OssConfigRemove(req *types.OssConfigRemoveReq) (resp *types.BaseResp, err error) {
	if req.OssConfigIds == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置ID不能为空",
		}, nil
	}

	// 1. 解析配置ID列表（逗号分隔）
	ossConfigIdStrs := strings.Split(req.OssConfigIds, ",")
	var ossConfigIds []int64
	for _, idStr := range ossConfigIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("配置ID格式错误: %s", idStr),
			}, nil
		}
		ossConfigIds = append(ossConfigIds, id)
	}

	if len(ossConfigIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置ID不能为空",
		}, nil
	}

	// 2. 删除配置
	for _, ossConfigId := range ossConfigIds {
		err = l.svcCtx.SysOssConfigModel.Delete(l.ctx, ossConfigId)
		if err != nil {
			if err == model.ErrNotFound {
				// 记录不存在，继续删除其他
				continue
			}
			l.Errorf("删除对象存储配置失败: ossConfigId=%d, err=%v", ossConfigId, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除对象存储配置失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
