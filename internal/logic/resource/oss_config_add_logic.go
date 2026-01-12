// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"database/sql"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssConfigAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 新增对象存储配置
func NewOssConfigAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssConfigAddLogic {
	return &OssConfigAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssConfigAddLogic) OssConfigAdd(req *types.OssConfigReq) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if req.ConfigKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "配置key不能为空",
		}, nil
	}
	if req.AccessKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "accessKey不能为空",
		}, nil
	}
	if req.SecretKey == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "秘钥不能为空",
		}, nil
	}
	if req.BucketName == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "桶名称不能为空",
		}, nil
	}
	if req.Endpoint == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "访问站点不能为空",
		}, nil
	}

	// 2. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 3. 生成主键ID（使用雪花算法）
	ossConfigId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成OSS配置ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成OSS配置ID失败",
		}, err
	}

	// 4. 构建OSS配置实体
	ossConfig := &model.SysOssConfig{
		OssConfigId:  ossConfigId,
		TenantId:     tenantId,
		ConfigKey:    req.ConfigKey,
		AccessKey:    req.AccessKey,
		SecretKey:    req.SecretKey,
		BucketName:   req.BucketName,
		Prefix:       req.Prefix,
		Endpoint:     req.Endpoint,
		Domain:       req.Domain,
		IsHttps:      req.IsHttps,
		Region:       req.Region,
		AccessPolicy: req.AccessPolicy,
		Status:       req.Status,
		Ext1:         req.Ext1,
		CreateDept:   sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		CreateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		Remark:       sql.NullString{String: req.Remark, Valid: req.Remark != ""},
	}

	// 5. 插入数据库
	_, err = l.svcCtx.SysOssConfigModel.Insert(l.ctx, ossConfig)
	if err != nil {
		l.Errorf("新增对象存储配置失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "新增对象存储配置失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
