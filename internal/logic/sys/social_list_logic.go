// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SocialListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询社会化关系列表
func NewSocialListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SocialListLogic {
	return &SocialListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SocialListLogic) SocialList() (resp *types.SocialListResp, err error) {
	// 从 context 中获取当前用户ID
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		return &types.SocialListResp{
			BaseResp: types.BaseResp{
				Code: 401,
				Msg:  "未授权，请先登录",
			},
			Data: []types.SysSocialVo{},
		}, nil
	}

	// 查询当前用户的社会化关系列表
	rows, err := l.svcCtx.SysSocialModel.FindByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询社会化关系列表失败: %v", err)
		return &types.SocialListResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询社会化关系列表失败",
			},
			Data: []types.SysSocialVo{},
		}, err
	}

	// 转换为 VO
	voList := make([]types.SysSocialVo, 0, len(rows))
	for _, row := range rows {
		vo := types.SysSocialVo{
			Id:          row.Id,
			UserId:      row.UserId,
			TenantId:    row.TenantId,
			AuthId:      row.AuthId,
			Source:      row.Source,
			UserName:    row.UserName,
			NickName:    row.NickName,
			Email:       row.Email,
			Avatar:      row.Avatar,
			AccessToken: row.AccessToken,
		}

		// 处理可空字段
		if row.OpenId.Valid {
			vo.OpenId = row.OpenId.String
		}
		if row.ExpireIn.Valid {
			vo.ExpireIn = row.ExpireIn.Int64
		}
		if row.RefreshToken.Valid {
			vo.RefreshToken = row.RefreshToken.String
		}
		if row.AccessCode.Valid {
			vo.AccessCode = row.AccessCode.String
		}
		if row.UnionId.Valid {
			vo.UnionId = row.UnionId.String
		}
		if row.Scope.Valid {
			vo.Scope = row.Scope.String
		}
		if row.TokenType.Valid {
			vo.TokenType = row.TokenType.String
		}
		if row.IdToken.Valid {
			vo.IdToken = row.IdToken.String
		}
		if row.MacAlgorithm.Valid {
			vo.MacAlgorithm = row.MacAlgorithm.String
		}
		if row.MacKey.Valid {
			vo.MacKey = row.MacKey.String
		}
		if row.Code.Valid {
			vo.Code = row.Code.String
		}
		if row.OauthToken.Valid {
			vo.OauthToken = row.OauthToken.String
		}
		if row.OauthTokenSecret.Valid {
			vo.OauthTokenSecret = row.OauthTokenSecret.String
		}
		if row.CreateTime.Valid {
			vo.CreateTime = row.CreateTime.Time.Format("2006-01-02 15:04:05")
		}

		voList = append(voList, vo)
	}

	return &types.SocialListResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "查询成功",
		},
		Data: voList,
	}, nil
}
