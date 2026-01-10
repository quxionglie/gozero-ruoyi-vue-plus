package sys

import (
	"context"
	"database/sql"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增用户
func NewUserAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAddLogic {
	return &UserAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAddLogic) UserAdd(req *types.UserReq) (resp *types.BaseResp, err error) {
	// 1. 校验用户名是否唯一
	unique, err := l.svcCtx.SysUserModel.CheckUserNameUnique(l.ctx, req.UserName, 0)
	if err != nil {
		l.Errorf("校验用户名唯一性失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "校验用户名唯一性失败",
		}, err
	}
	if !unique {
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增用户'" + req.UserName + "'失败，登录账号已存在",
		}, nil
	}

	// 2. 校验手机号是否唯一
	if req.Phonenumber != "" {
		unique, err := l.svcCtx.SysUserModel.CheckPhoneUnique(l.ctx, req.Phonenumber, 0)
		if err != nil {
			l.Errorf("校验手机号唯一性失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "校验手机号唯一性失败",
			}, err
		}
		if !unique {
			return &types.BaseResp{
				Code: 500,
				Msg:  "新增用户'" + req.UserName + "'失败，手机号码已存在",
			}, nil
		}
	}

	// 3. 校验邮箱是否唯一
	if req.Email != "" {
		unique, err := l.svcCtx.SysUserModel.CheckEmailUnique(l.ctx, req.Email, 0)
		if err != nil {
			l.Errorf("校验邮箱唯一性失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "校验邮箱唯一性失败",
			}, err
		}
		if !unique {
			return &types.BaseResp{
				Code: 500,
				Msg:  "新增用户'" + req.UserName + "'失败，邮箱账号已存在",
			}, nil
		}
	}

	// 4. 获取当前用户信息
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取用户ID失败",
		}, err
	}

	tenantId, err := util.GetTenantIdFromContext(l.ctx)
	if err != nil {
		l.Errorf("获取租户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取租户ID失败",
		}, err
	}

	// 5. 查询当前用户的部门ID
	currentUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询当前用户信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询当前用户信息失败",
		}, err
	}
	var createDept int64 = 0
	if currentUser.CreateDept.Valid {
		createDept = currentUser.CreateDept.Int64
	} else if currentUser.DeptId.Valid {
		createDept = currentUser.DeptId.Int64
	}

	// 6. 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "密码加密失败",
		}, err
	}

	// 7. 生成用户ID
	newUserId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成用户ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成用户ID失败",
		}, err
	}

	// 8. 构建用户实体
	newUser := &model.SysUser{
		UserId:      newUserId,
		TenantId:    tenantId,
		UserName:    req.UserName,
		NickName:    req.NickName,
		UserType:    "00", // sys_user系统用户
		Email:       req.Email,
		Phonenumber: req.Phonenumber,
		Sex:         req.Sex,
		Password:    string(hashedPassword),
		Status:      req.Status,
		DelFlag:     "0",
		CreateBy:    sql.NullInt64{Int64: userId, Valid: true},
		CreateDept:  sql.NullInt64{Int64: createDept, Valid: createDept > 0},
	}
	if req.DeptId > 0 {
		newUser.DeptId = sql.NullInt64{Int64: req.DeptId, Valid: true}
	}
	// Avatar 字段在 UserReq 中是可选的，如果存在则设置
	if req.Remark != "" {
		newUser.Remark = sql.NullString{String: req.Remark, Valid: true}
	}

	// 9. 插入用户
	_, err = l.svcCtx.SysUserModel.Insert(l.ctx, newUser)
	if err != nil {
		l.Errorf("新增用户失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增用户失败",
		}, err
	}

	// 10. 插入用户角色关联
	if len(req.RoleIds) > 0 {
		// 非超级管理员，禁止包含超级管理员角色
		isSuperAdmin, err := l.svcCtx.SysRoleModel.CheckIsSuperAdmin(l.ctx, userId)
		if err != nil {
			l.Errorf("检查超级管理员失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "检查超级管理员失败",
			}, err
		}

		var roleIds []int64
		for _, roleId := range req.RoleIds {
			// 非超级管理员，禁止包含超级管理员角色（roleId=1）
			if !isSuperAdmin && roleId == 1 {
				continue
			}
			roleIds = append(roleIds, roleId)
		}

		if len(roleIds) > 0 {
			err = l.svcCtx.SysUserRoleModel.InsertBatchByUserId(l.ctx, newUserId, roleIds)
			if err != nil {
				l.Errorf("新增用户角色关联失败: %v", err)
				return &types.BaseResp{
					Code: 500,
					Msg:  "新增用户角色关联失败",
				}, err
			}
		}
	}

	// 11. 插入用户岗位关联
	if len(req.PostIds) > 0 {
		err = l.svcCtx.SysUserPostModel.InsertBatch(l.ctx, newUserId, req.PostIds)
		if err != nil {
			l.Errorf("新增用户岗位关联失败: %v", err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "新增用户岗位关联失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "新增成功",
	}, nil
}
