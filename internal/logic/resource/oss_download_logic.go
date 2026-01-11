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
	oss, err := l.svcCtx.SysOssModel.FindOne(l.ctx, req.OssId)
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
	originalName := oss.OriginalName
	// URL编码文件名，避免中文乱码
	encodedName := url.QueryEscape(originalName)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName))
	w.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")

	// 3. TODO: 从OSS服务下载文件
	// 实际应该调用OSS服务下载文件并写入响应
	// storage := OssFactory.instance(oss.Service)
	// storage.download(oss.FileName, w)

	// 临时方案：返回一个错误提示（因为需要实际的OSS服务实现）
	http.Error(w, "OSS服务未配置，无法下载文件", http.StatusNotImplemented)
	return fmt.Errorf("OSS服务未配置")
}
