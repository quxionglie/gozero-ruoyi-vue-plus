package sys

import (
	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/types"
)

// convertUserToVo 转换用户实体为响应格式
func convertUserToVo(user *model.SysUser) types.SysUserVo {
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
		Avatar:      0,
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
		// TODO: 查询部门名称
	}
	if user.Avatar.Valid {
		userVo.Avatar = user.Avatar.Int64
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
