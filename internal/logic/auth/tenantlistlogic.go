// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"database/sql"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TenantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 登录页面租户下拉框，获取可用租户列表
func NewTenantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TenantListLogic {
	return &TenantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TenantListLogic) TenantList() (resp *types.TenantListResp, err error) {
	// 检查租户是否启用
	tenantEnabled := l.isTenantEnabled()

	resp = &types.TenantListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.LoginTenantVo{
			TenantEnabled: tenantEnabled,
		},
	}

	// 如果租户未启用，直接返回
	if !tenantEnabled {
		return resp, nil
	}

	// 查询所有可用的租户
	tenants, err := l.svcCtx.SysTenantModel.FindAllAvailable(l.ctx)
	if err != nil {
		l.Errorf("查询租户列表失败: %v", err)
		return &types.TenantListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询租户列表失败",
			},
		}, err
	}

	// 转换为响应格式
	voList := make([]types.TenantListVo, 0, len(tenants))
	for _, tenant := range tenants {
		companyName := ""
		if tenant.CompanyName.Valid {
			companyName = tenant.CompanyName.String
		}
		domain := ""
		if tenant.Domain.Valid {
			domain = tenant.Domain.String
		}

		voList = append(voList, types.TenantListVo{
			TenantId:    tenant.TenantId,
			CompanyName: companyName,
			Domain:      domain,
		})
	}

	resp.Data.VoList = voList

	return resp, nil
}

// isTenantEnabled 检查租户是否启用
func (l *TenantListLogic) isTenantEnabled() bool {
	// 查询系统配置：sys.tenant.enable
	query := "SELECT config_value FROM sys_config WHERE config_key = 'sys.tenant.enable' AND del_flag = '0' LIMIT 1"
	var configValue string
	err := l.svcCtx.DB.QueryRowCtx(l.ctx, &configValue, query)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果配置不存在，默认启用租户
			return true
		}
		l.Errorf("查询租户配置失败: %v", err)
		// 出错时默认启用租户
		return true
	}

	// 配置值为 "true" 或 "1" 表示启用
	return configValue == "true" || configValue == "1"
}
