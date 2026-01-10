// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	model "gozero-ruoyi-vue-plus/internal/model/sys"
	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除字典类型
func NewDictTypeRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeRemoveLogic {
	return &DictTypeRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeRemoveLogic) DictTypeRemove(req *types.DictTypeRemoveReq) (resp *types.BaseResp, err error) {
	if req.DictIds == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典ID不能为空",
		}, nil
	}

	// 解析字典ID列表（逗号分隔）
	dictIdStrs := strings.Split(req.DictIds, ",")
	var dictIds []int64
	for _, idStr := range dictIdStrs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("字典ID格式错误: %s", idStr),
			}, nil
		}
		dictIds = append(dictIds, id)
	}

	if len(dictIds) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典ID不能为空",
		}, nil
	}

	// 检查是否已分配字典数据
	for _, dictId := range dictIds {
		dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, dictId)
		if err != nil {
			if err != model.ErrNotFound {
				l.Errorf("查询字典类型失败: %v", err)
			}
			continue
		}

		// 检查是否有字典数据使用该类型
		count, err := l.svcCtx.SysDictDataModel.CountByDictType(l.ctx, dictType.DictType)
		if err == nil && count > 0 {
			return &types.BaseResp{
				Code: 500,
				Msg:  fmt.Sprintf("%s已分配,不能删除", dictType.DictName),
			}, nil
		}
	}

	// 删除字典类型
	for _, dictId := range dictIds {
		err = l.svcCtx.SysDictTypeModel.Delete(l.ctx, dictId)
		if err != nil {
			l.Errorf("删除字典类型失败: dictId=%d, err=%v", dictId, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除字典类型失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
