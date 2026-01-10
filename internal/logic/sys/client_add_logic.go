// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增客户端管理
func NewClientAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientAddLogic {
	return &ClientAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientAddLogic) ClientAdd(req *types.ClientReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
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
	unique, err := l.svcCtx.SysClientModel.CheckClientKeyUnique(l.ctx, req.ClientKey, 0)
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
			Msg:  fmt.Sprintf("新增客户端'%s'失败，客户端key已存在", req.ClientKey),
		}, nil
	}

	// 3. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 4. 生成主键ID（使用雪花算法）
	id, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成客户端ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成客户端ID失败",
		}, err
	}

	// 5. 生成 clientId（MD5(clientKey + clientSecret)）
	clientKeySecret := req.ClientKey + req.ClientSecret
	hash := md5.Sum([]byte(clientKeySecret))
	clientId := fmt.Sprintf("%x", hash)

	// 6. 将 grantTypeList 转换为逗号分隔的字符串
	grantTypeStr := strings.Join(req.GrantTypeList, ",")

	// 7. 构建客户端实体
	client := &model.SysClient{
		Id:            id,
		ClientId:      sql.NullString{String: clientId, Valid: true},
		ClientKey:     sql.NullString{String: req.ClientKey, Valid: true},
		ClientSecret:  sql.NullString{String: req.ClientSecret, Valid: true},
		GrantType:     sql.NullString{String: grantTypeStr, Valid: true},
		DeviceType:    sql.NullString{String: req.DeviceType, Valid: req.DeviceType != ""},
		ActiveTimeout: req.ActiveTimeout,
		Timeout:       req.Timeout,
		Status:        req.Status,
		DelFlag:       "0",
		CreateDept:    sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:      sql.NullInt64{Int64: userId, Valid: userId > 0},
	}
	if client.Status == "" {
		client.Status = "0"
	}
	if client.ActiveTimeout == 0 {
		client.ActiveTimeout = 0
	}
	if client.Timeout == 0 {
		client.Timeout = 0
	}

	// 8. 插入数据库
	_, err = l.svcCtx.SysClientModel.Insert(l.ctx, client)
	if err != nil {
		l.Errorf("新增客户端管理失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增客户端管理失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
