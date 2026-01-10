// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package monitor

import (
	"context"
	"encoding/json"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OnlineListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取在线用户监控列表
func NewOnlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineListLogic {
	return &OnlineListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineListLogic) OnlineList(req *types.OnlineListReq) (resp *types.OnlineListResp, err error) {
	// 获取所有 online_tokens:* 的 keys
	keys, err := l.svcCtx.RedisConn.KeysCtx(l.ctx, "online_tokens:*")
	if err != nil {
		l.Errorf("获取在线token列表失败: %v", err)
		return &types.OnlineListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "获取在线用户列表失败",
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

		// 过滤条件
		if req.Ipaddr != "" {
			if ipaddr, ok := onlineInfo["ipaddr"].(string); !ok || ipaddr != req.Ipaddr {
				continue
			}
		}
		if req.UserName != "" {
			if userName, ok := onlineInfo["userName"].(string); !ok || userName != req.UserName {
				continue
			}
		}

		vo := convertOnlineInfoToVo(onlineInfo, token)
		if vo != nil {
			onlineUserList = append(onlineUserList, *vo)
		}
	}

	// 反转列表，让最新的在线用户在前面
	for i, j := 0, len(onlineUserList)-1; i < j; i, j = i+1, j-1 {
		onlineUserList[i], onlineUserList[j] = onlineUserList[j], onlineUserList[i]
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
