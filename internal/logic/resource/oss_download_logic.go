// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package resource

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"
	"gozero-ruoyi-vue-plus/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type OssDownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 下载OSS对象
func NewOssDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OssDownloadLogic {
	return &OssDownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OssDownloadLogic) OssDownload(req *types.OssDownloadReq, w http.ResponseWriter, r *http.Request) error {
	// 1. 查询OSS对象信息
	ossObj, err := l.svcCtx.SysOssModel.FindOne(l.ctx, req.OssId)
	if err != nil {
		if err == model.ErrNotFound {
			http.Error(w, "文件数据不存在", http.StatusNotFound)
			return nil
		}
		l.Errorf("查询OSS对象信息失败: %v", err)
		http.Error(w, "查询OSS对象信息失败", http.StatusInternalServerError)
		return err
	}

	// 2. 设置响应头
	originalName := ossObj.OriginalName
	// URL编码文件名，避免中文乱码
	encodedName := url.QueryEscape(originalName)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	w.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")

	// 3. 获取OSS客户端并下载文件
	tenantId, _ := util.GetTenantIdFromContext(l.ctx)
	if ossObj.Service == "" {
		http.Error(w, "OSS服务商未配置", http.StatusBadRequest)
		return fmt.Errorf("OSS服务商未配置")
	}
	ossClient, err := l.svcCtx.OssManager.GetClientByConfigKey(l.ctx, ossObj.Service, tenantId)
	if err != nil {
		l.Errorf("获取OSS客户端失败: %v", err)
		http.Error(w, "获取OSS客户端失败: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	// 下载文件
	err = ossClient.Download(ossObj.FileName, w)
	if err != nil {
		l.Errorf("从OSS下载文件失败: %v", err)
		http.Error(w, "下载文件失败: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
