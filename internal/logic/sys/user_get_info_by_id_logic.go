package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserGetInfoByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户编号获取详细信息
func NewUserGetInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetInfoByIdLogic {
	return &UserGetInfoByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserGetInfoByIdLogic) UserGetInfoById(req *types.UserGetInfoReq) (resp *types.UserDetailResp, err error) {
	// 1. 查询用户信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		if err == model.ErrNotFound {
			return &types.UserDetailResp{
				BaseResp: types.BaseResp{
					Code: 404,
					Msg:  "用户不存在",
				},
			}, nil
		}
		l.Errorf("查询用户信息失败: %v", err)
		return &types.UserDetailResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户信息失败",
			},
		}, err
	}

	// 2. 查询用户角色
	roles, err := l.svcCtx.SysRoleModel.SelectRolesByUserId(l.ctx, req.UserId)
	if err != nil {
		l.Errorf("查询用户角色失败: %v", err)
		return &types.UserDetailResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "查询用户角色失败",
			},
		}, err
	}

	// 3. 转换为响应格式
	userVo := convertUserToVo(l.ctx, l.svcCtx, user)
	roleVos := make([]types.SysRoleVo, 0, len(roles))
	for _, role := range roles {
		roleVo := types.SysRoleVo{
			RoleId:            role.RoleId,
			RoleName:          role.RoleName,
			RoleKey:           role.RoleKey,
			RoleSort:          int32(role.RoleSort),
			DataScope:         role.DataScope,
			MenuCheckStrictly: nil, // 匹配 Java 返回 null
			DeptCheckStrictly: nil, // 匹配 Java 返回 null
			Status:            role.Status,
			Remark:            nil,
			CreateTime:        nil,
			SuperAdmin:        role.RoleId == 1 || role.RoleKey == "superadmin",
			Flag:              false, // 默认值
		}
		// Remark 现在是 *string，当值为 null 时返回 nil
		if role.Remark.Valid {
			remarkStr := role.Remark.String
			roleVo.Remark = &remarkStr
		}
		// CreateTime 现在是 *string，当值为 null 时返回 nil
		if role.CreateTime.Valid {
			createTimeStr := role.CreateTime.Time.Format("2006-01-02 15:04:05")
			roleVo.CreateTime = &createTimeStr
		}
		roleVos = append(roleVos, roleVo)
	}
	userVo.Roles = roleVos

	return &types.UserDetailResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.UserInfoDataVo{
			User: userVo,
		},
	}, nil
}
