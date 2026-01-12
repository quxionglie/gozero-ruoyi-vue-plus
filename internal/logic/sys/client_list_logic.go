// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"strings"

	"gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询客户端管理列表
func NewClientListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientListLogic {
	return &ClientListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientListLogic) ClientList(req *types.ClientListReq) (resp *types.TableDataInfoResp, err error) {
	// 1. 构建查询条件
	clientQuery := &sys.ClientQuery{
		ClientId:     req.ClientId,
		ClientKey:    req.ClientKey,
		ClientSecret: req.ClientSecret,
		Status:       req.Status,
	}
	pageQuery := &sys.PageQuery{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		OrderByColumn: req.OrderByColumn,
		IsAsc:         req.IsAsc,
	}

	// 3. 查询数据
	clients, total, err := l.svcCtx.SysClientModel.FindPage(l.ctx, clientQuery, pageQuery)
	if err != nil {
		l.Errorf("查询客户端管理列表失败: %v", err)
		return &types.TableDataInfoResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询客户端管理列表失败",
			},
		}, err
	}

	// 4. 转换为响应格式
	rows := make([]types.ClientVo, 0, len(clients))
	for _, client := range clients {
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
		rows = append(rows, clientVo)
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
