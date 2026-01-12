package sys

import (
	"context"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
)

// convertUserToVo 转换用户实体为响应格式
func convertUserToVo(ctx context.Context, svcCtx *svc.ServiceContext, user *model.SysUser) types.SysUserVo {
	userVo := types.SysUserVo{
		UserId:      user.UserId,
		TenantId:    user.TenantId,
		DeptId:      0,
		UserName:    user.UserName,
		NickName:    user.NickName,
		UserType:    user.UserType,
		Email:       user.Email,
		Phonenumber: user.Phonenumber,
		Sex:         user.Sex,
		Avatar:      nil, // *string 类型，初始化为 nil
		Status:      user.Status,
		LoginIp:     user.LoginIp,
		LoginDate:   "",
		Remark:      "",
		CreateTime:  "",
		DeptName:    "",
		Roles:       []types.SysRoleVo{},
	}

	if user.DeptId.Valid {
		userVo.DeptId = user.DeptId.Int64
		// 查询部门名称
		dept, err := svcCtx.SysDeptModel.FindOne(ctx, user.DeptId.Int64)
		if err == nil {
			userVo.DeptName = dept.DeptName
		}
	}
	// Avatar 现在是 *string，从 sys_oss 表查询 URL
	if user.Avatar.Valid {
		oss, err := svcCtx.SysOssModel.FindOne(ctx, user.Avatar.Int64)
		if err == nil {
			userVo.Avatar = &oss.Url
		} else {
			// 如果查询失败，返回 nil（可能 OSS 记录已删除）
			userVo.Avatar = nil
		}
	} else {
		userVo.Avatar = nil
	}
	if user.LoginDate.Valid {
		userVo.LoginDate = user.LoginDate.Time.Format("2006-01-02 15:04:05")
	}
	if user.Remark.Valid {
		userVo.Remark = user.Remark.String
	}
	if user.CreateTime.Valid {
		userVo.CreateTime = user.CreateTime.Time.Format("2006-01-02 15:04:05")
	}

	return userVo
}
