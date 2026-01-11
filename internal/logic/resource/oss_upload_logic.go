// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传OSS对象存储
func NewOssUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssUploadLogic {
	return &OssUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssUploadLogic) OssUpload(file multipart.File, fileHeader *multipart.FileHeader) (resp *types.OssUploadResp, err error) {
	// 1. 参数校验
	if file == nil || fileHeader == nil {
		return &types.OssUploadResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "上传文件不能为空",
			},
		}, nil
	}

	// 2. 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		l.Errorf("读取文件内容失败: %v", err)
		return &types.OssUploadResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "读取文件内容失败",
			},
		}, err
	}

	if len(fileBytes) == 0 {
		return &types.OssUploadResp{
			BaseResp: types.BaseResp{
				Code: 400,
				Msg:  "上传文件不能为空",
			},
		}, nil
	}

	// 3. 获取文件信息
	originalFileName := fileHeader.Filename
	fileSuffix := filepath.Ext(originalFileName)
	if fileSuffix == "" {
		fileSuffix = ""
	}

	// 4. 获取当前用户信息
	userId, _ := util.GetUserIdFromContext(l.ctx)
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)

	// 获取用户的部门ID
	var deptId int64
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err == nil && user.DeptId.Valid {
		deptId = user.DeptId.Int64
	}

	// 5. 生成主键ID（使用雪花算法）
	newOssId, err := util.GenerateID()
	if err != nil {
		l.Errorf("生成OSS ID失败: %v", err)
		return &types.OssUploadResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "生成OSS ID失败",
			},
		}, err
	}

	// 6. TODO: 上传文件到OSS服务（这里简化处理，直接保存URL）
	// 实际应该调用OSS服务上传文件，获取文件URL
	// storage := OssFactory.instance()
	// uploadResult := storage.uploadSuffix(fileBytes, fileSuffix, fileHeader.Header.Get("Content-Type"))
	// fileUrl := uploadResult.getUrl()
	// fileName := uploadResult.getFilename()

	// 临时方案：生成一个模拟的URL和文件名
	fileUrl := fmt.Sprintf("/oss/%d%s", newOssId, fileSuffix)
	fileName := fmt.Sprintf("%d%s", newOssId, fileSuffix)

	// 7. 构建OSS实体
	oss := &model.SysOss{
		OssId:        newOssId,
		TenantId:     tenantId,
		FileName:     fileName,
		OriginalName: originalFileName,
		FileSuffix:   fileSuffix,
		Url:          fileUrl,
		Ext1:         sql.NullString{String: fmt.Sprintf(`{"fileSize":%d,"contentType":"%s"}`, len(fileBytes), fileHeader.Header.Get("Content-Type")), Valid: true},
		Service:      "local", // TODO: 从配置中获取默认服务商
		CreateDept:   sql.NullInt64{Int64: deptId, Valid: deptId > 0},
		CreateBy:     sql.NullInt64{Int64: userId, Valid: userId > 0},
	}

	// 8. 插入数据库
	_, err = l.svcCtx.SysOssModel.Insert(l.ctx, oss)
	if err != nil {
		l.Errorf("保存OSS对象失败: %v", err)
		return &types.OssUploadResp{
			BaseResp: types.BaseResp{
				Code: 500,
				Msg:  "保存OSS对象失败",
			},
		}, err
	}

	return &types.OssUploadResp{
		BaseResp: types.BaseResp{
			Code: 200,
			Msg:  "操作成功",
		},
		Data: types.OssUploadVo{
			Url:      fileUrl,
			FileName: fileName,
			OssId:    fmt.Sprintf("%d", newOssId),
		},
	}, nil
}
