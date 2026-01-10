// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询客户端管理详细
func NewClientGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientGetInfoLogic {
	return &ClientGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientGetInfoLogic) ClientGetInfo(req *types.ClientGetInfoReq) (resp *types.ClientResp, err error) {
	// 1. 查询客户端信息
	client, err := l.svcCtx.SysClientModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.ClientResp{
				BaseResp: types.BaseResp{
					Code: 500,
					Msg:  "客户端不存在",
				},
			}, nil
		}
		l.Errorf("查询客户端信息失败: %v", err)
		return &types.ClientResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询客户端信息失败",
			},
		}, err
	}

	// 2. 转换为响应格式
	clientVo := types.ClientVo{
		Id:            client.Id,
		ClientId:      "",
		ClientKey:     "",
		ClientSecret:  "",
		GrantTypeList: []string{},
		GrantType:     "",
		DeviceType:    "",
		ActiveTimeout: client.ActiveTimeout,
		Timeout:       client.Timeout,
		Status:        client.Status,
	}
	if client.ClientId.Valid {
		clientVo.ClientId = client.ClientId.String
	}
	if client.ClientKey.Valid {
		clientVo.ClientKey = client.ClientKey.String
	}
	if client.ClientSecret.Valid {
		clientVo.ClientSecret = client.ClientSecret.String
	}
	if client.GrantType.Valid && client.GrantType.String != "" {
		// 将逗号分隔的字符串转换为数组
		clientVo.GrantType = client.GrantType.String
		clientVo.GrantTypeList = strings.Split(client.GrantType.String, ",")
		// 去除空格
		for i, v := range clientVo.GrantTypeList {
			clientVo.GrantTypeList[i] = strings.TrimSpace(v)
		}
	}
	if client.DeviceType.Valid {
		clientVo.DeviceType = client.DeviceType.String
	}

	return &types.ClientResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: clientVo,
	}, nil
}
