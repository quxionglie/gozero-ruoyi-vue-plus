package sys

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 头像上传
func NewUserProfileAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileAvatarLogic {
	return &UserProfileAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileAvatarLogic) UserProfileAvatar(file multipart.File, fileHeader *multipart.FileHeader) (resp *types.BaseResp, err error) {
	// 1. 参数校验
	if file == nil || fileHeader == nil {
		return &types.BaseResp{
			Code: 400,
			Msg:  "上传文件不能为空",
		}, nil
	}

	// 2. 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		l.Errorf("读取文件内容失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "读取文件内容失败",
		}, err
	}

	if len(fileBytes) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "上传文件不能为空",
		}, nil
	}

	// 3. 获取文件信息
	originalFileName := fileHeader.Filename
	fileSuffix := filepath.Ext(originalFileName)
	if fileSuffix == "" {
		fileSuffix = ".jpg" // 默认jpg
	}

	// 4. 获取当前用户信息
	userId, err := util.GetUserIdFromContext(l.ctx)
	if err != nil {
		return &types.BaseResp{
			Code: 401,
			Msg:  "未授权，请先登录",
		}, nil
	}

	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户信息失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "查询用户信息失败",
		}, err
	}
	if user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 5. 生成主键ID（使用雪花算法）
	newOssId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成OSS ID失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "生成OSS ID失败",
		}, err
	}

	// 6. 获取OSS客户端并上传文件
	ossClient, err := l.svcCtx.OssManager.GetDefaultClient(l.ctx, tenantId)
	if err != nil {
		l.Errorf("获取OSS客户端失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "获取OSS客户端失败: " + err.Error(),
		}, err
	}

	// 获取Content-Type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // 默认图片类型
	}

	// 上传文件
	uploadResult, err := ossClient.UploadSuffix(fileBytes, fileSuffix, contentType)
	if err != nil {
		l.Errorf("上传文件到OSS失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "上传文件到OSS失败: " + err.Error(),
		}, err
	}

	fileUrl := uploadResult.URL
	fileName := uploadResult.Filename
	configKey := ossClient.GetConfigKey()

	// 7. 构建OSS实体
	oss := &model.SysOss{
		OssId:        newOssId,
		TenantId:     tenantId,
		FileName:     fileName,
		OriginalName: originalFileName,
		FileSuffix:   fileSuffix,
		Url:          fileUrl,
		Ext1:         sql.NullString{String: fmt.Sprintf(`{"fileSize":%d,"contentType":"%s"}`, len(fileBytes), contentType), Valid: true},
		Service:      configKey,
		CreateDept:   sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
		CreateTime:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdateTime:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	// 8. 插入数据库
	_, err = l.svcCtx.SysOssModel.Insert(l.ctx, oss)
	if err != nil {
		l.Errorf("保存OSS对象失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "保存OSS对象失败",
		}, err
	}

	// 9. 更新用户头像字段（使用OSS ID作为头像标识）
	err = l.svcCtx.SysUserModel.UpdateUserAvatar(l.ctx, userId, newOssId)
	if err != nil {
		l.Errorf("更新用户头像失败: %v", err)
		return &types.BaseResp{
			Code: 500,
			Msg:  "更新用户头像失败",
		}, err
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "上传成功",
	}, nil
}
