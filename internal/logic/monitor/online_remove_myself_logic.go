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

type OnlineRemoveMyselfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 强退当前在线设备
func NewOnlineRemoveMyselfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OnlineRemoveMyselfLogic {
	return &OnlineRemoveMyselfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OnlineRemoveMyselfLogic) OnlineRemoveMyself(req *types.OnlineRemoveReq) (resp *types.BaseResp, err error) {
	// 获取当前用户名
	username, err := util.GetUsernameFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户名失败: %v", err)
		return &types.BaseResp{
			Code: 401,
			Msg:  "未登录",
		}, err
	}

	// 检查该 token 是否属于当前用户
	onlineTokenKey := "online_tokens:" + req.TokenId
	infoJSON, getErr := l.svcCtx.RedisConn.GetCtx(l.ctx, onlineTokenKey)
	if getErr != nil {
		return &types.BaseResp{
			Code: 404,
			Msg:  "未找到该在线设备",
		}, getErr
	}

	var onlineInfo map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(infoJSON), &onlineInfo)
	if unmarshalErr != nil {
		return &types.BaseResp{
			Code: 500,
			Msg:  "解析在线设备信息失败",
		}, unmarshalErr
	}

	// 验证该 token 是否属于当前用户
	if userName, ok := onlineInfo["userName"].(string); !ok || userName != username {
		return &types.BaseResp{
			Code: 403,
			Msg:  "无权操作该设备",
		}, nil
	}

	// 删除在线 token 缓存
	_, delErr := l.svcCtx.RedisConn.DelCtx(l.ctx, onlineTokenKey)
	if delErr != nil {
		l.Errorf("强退当前在线设备失败: %v", delErr)
		return &types.BaseResp{
			Code: 500,
			Msg:  "强退当前在线设备失败",
		}, delErr
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "强退成功",
	}, nil
}
