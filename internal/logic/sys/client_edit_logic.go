// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改客户端管理
func NewClientEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientEditLogic {
	return &ClientEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientEditLogic) ClientEdit(req *types.ClientReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.Id <= 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "客户端ID不能为空",
		}, nil
	}
	if req.ClientKey == "" {
		return &types.BaseResp{Code: 400, Msg: "客户端key不能为空"}, nil
	}
	if req.ClientSecret == "" {
		return &types.BaseResp{Code: 400, Msg: "客户端秘钥不能为空"}, nil
	}
	if len(req.GrantTypeList) == 0 {
		return &types.BaseResp{Code: 400, Msg: "授权类型不能为空"}, nil
	}

	// 2. 校验客户端key唯一性
	unique, err := l.svcCtx.SysClientModel.CheckClientKeyUnique(l.ctx, req.ClientKey, req.Id)
	if err != nil {
		l.Errorf("校验客户端key唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验客户端key唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  fmt.Sprintf("修改客户端'%s'失败，客户端key已存在", req.ClientKey),
		}, nil
	}

	// 3. 检查客户端是否存在
	_, err = l.svcCtx.SysClientModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.BaseResp{
				Code: 500,
				Msg:  "客户端不存在",
			}, nil
		}
		l.Errorf("查询客户端失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询客户端失败",
		}, err
	}

	// 4. 获取当前用户ID
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 5. 生成新的 clientId（MD5(clientKey + clientSecret)）
	clientKeySecret := req.ClientKey + req.ClientSecret
	hash := md5.Sum([]byte(clientKeySecret))
	clientId := fmt.Sprintf("%x", hash)

	// 6. 将 grantTypeList 转换为逗号分隔的字符串
	grantTypeStr := strings.Join(req.GrantTypeList, ",")

	// 7. 更新客户端信息（只设置表单输入的字段）
	updateClient := &model.SysClient{
		Id:           req.Id,
		ClientId:     sql.NullString{String: clientId, Valid: true},
		ClientKey:    sql.NullString{String: req.ClientKey, Valid: true},
		ClientSecret: sql.NullString{String: req.ClientSecret, Valid: true},
		GrantType:    sql.NullString{String: grantTypeStr, Valid: true},
		UpdateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	if req.DeviceType != "" {
		updateClient.DeviceType = sql.NullString{String: req.DeviceType, Valid: true}
	}
	if req.ActiveTimeout > 0 {
		updateClient.ActiveTimeout = req.ActiveTimeout
	}
	if req.Timeout > 0 {
		updateClient.Timeout = req.Timeout
	}
	if req.Status != "" {
		updateClient.Status = req.Status
	}

	// 8. 更新数据库
	_, err = l.svcCtx.SysClientModel.UpdateById(l.ctx, updateClient)
	if err != nil {
		l.Errorf("修改客户端管理失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "修改客户端管理失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
