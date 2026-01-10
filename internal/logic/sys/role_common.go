package sys

import (
	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/types"
)

// convertRoleToVo 通用函数：转换角色实体为响应格式
func convertRoleToVo(role *model.SysRole) types.RoleVo {
	roleVo := types.RoleVo{
		RoleId:            role.RoleId,
		RoleName:          role.RoleName,
		RoleKey:           role.RoleKey,
		RoleSort:          int32(role.RoleSort),
		DataScope:         role.DataScope,
		MenuCheckStrictly: role.MenuCheckStrictly > 0,
		DeptCheckStrictly: role.DeptCheckStrictly > 0,
		Status:            role.Status,
		DelFlag:           role.DelFlag,
		CreateBy:          0,
		CreateTime:        "",
		UpdateBy:          0,
		UpdateTime:        "",
		Remark:            "",
		MenuIds:           []int64{},
		DeptIds:           []int64{},
	}

	if role.CreateBy.Valid {
		roleVo.CreateBy = role.CreateBy.Int64
	}
	if role.CreateTime.Valid {
		roleVo.CreateTime = role.CreateTime.Time.Format("2006-01-02 15:04:05")
	}
	if role.UpdateBy.Valid {
		roleVo.UpdateBy = role.UpdateBy.Int64
	}
	if role.UpdateTime.Valid {
		roleVo.UpdateTime = role.UpdateTime.Time.Format("2006-01-02 15:04:05")
	}
	if role.Remark.Valid {
		roleVo.Remark = role.Remark.String
	}

	return roleVo
}
