// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"
	"encoding/json"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineGetInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户登录在线设备
func NewOnlineGetInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineGetInfoLogic {
	return &OnlineGetInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineGetInfoLogic) OnlineGetInfo() (resp *types.OnlineListResp, err error) {
	// 获取当前用户名
	username, err := util.GetUsernameFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户名失败: %v", err)
		return &types.OnlineListResp{
			BaseResp: types.BaseResp{
				Code: 401,
				Msg:  "未登录",
			},
		}, err
	}

	// 获取所有 online_tokens:* 的 keys
	keys, err := l.svcCtx.RedisConn.KeysCtx(l.ctx, "online_tokens:*")
	if err != nil {
		l.Errorf("获取在线token列表失败: %v", err)
		return &types.OnlineListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取在线设备失败",
			},
		}, err
	}

	onlineUserList := make([]types.SysUserOnlineVo, 0)

	// 遍历 keys，获取在线用户信息
	for _, key := range keys {
		// 提取 token
		token := key[len("online_tokens:"):]

		// 获取在线用户信息
		infoJSON, getErr := l.svcCtx.RedisConn.GetCtx(l.ctx, key)
		if getErr != nil {
			continue
		}

		var onlineInfo map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(infoJSON), &onlineInfo)
		if unmarshalErr != nil {
			continue
		}

		// 检查是否属于当前用户
		if userName, ok := onlineInfo["userName"].(string); ok && userName == username {
			vo := convertOnlineInfoToVo(onlineInfo, token)
			if vo != nil {
				onlineUserList = append(onlineUserList, *vo)
			}
		}
	}

	return &types.OnlineListResp{
		Total: int64(len(onlineUserList)),
		Rows:  onlineUserList,
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
	}, nil
}

// convertOnlineInfoToVo 将在线用户信息转换为 VO
func convertOnlineInfoToVo(onlineInfo map[string]interface{}, tokenId string) *types.SysUserOnlineVo {
	vo := &types.SysUserOnlineVo{
		TokenId: tokenId,
	}

	if userName, ok := onlineInfo["userName"].(string); ok {
		vo.UserName = userName
	}
	if deptName, ok := onlineInfo["deptName"].(string); ok {
		vo.DeptName = deptName
	}
	if clientKey, ok := onlineInfo["clientKey"].(string); ok {
		vo.ClientKey = clientKey
	}
	if deviceType, ok := onlineInfo["deviceType"].(string); ok {
		vo.DeviceType = deviceType
	}
	if ipaddr, ok := onlineInfo["ipaddr"].(string); ok {
		vo.Ipaddr = ipaddr
	}
	if loginLocation, ok := onlineInfo["loginLocation"].(string); ok {
		vo.LoginLocation = loginLocation
	}
	if browser, ok := onlineInfo["browser"].(string); ok {
		vo.Browser = browser
	}
	if os, ok := onlineInfo["os"].(string); ok {
		vo.Os = os
	}
	if loginTime, ok := onlineInfo["loginTime"].(float64); ok {
		vo.LoginTime = int64(loginTime)
	} else if loginTime, ok := onlineInfo["loginTime"].(int64); ok {
		vo.LoginTime = loginTime
	}

	return vo
}
