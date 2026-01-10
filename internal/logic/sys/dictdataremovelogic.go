// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sys

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"gozero-ruoyi-vue-plus/internal/svc"
	"gozero-ruoyi-vue-plus/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictDataRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除字典数据
func NewDictDataRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictDataRemoveLogic {
	return &DictDataRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictDataRemoveLogic) DictDataRemove(req *types.DictDataRemoveReq) (resp *types.BaseResp, err error) {
	if req.DictCodes == "" {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典编码不能为空",
		}, nil
	}

	// 解析字典编码列表
	dictCodeStrs := strings.Split(req.DictCodes, ",")
	var dictCodes []int64
	for _, codeStr := range dictCodeStrs {
		codeStr = strings.TrimSpace(codeStr)
		if codeStr == "" {
			continue
		}
		code, err := strconv.ParseInt(codeStr, 10, 64)
		if err != nil {
			return &types.BaseResp{
				Code: 400,
				Msg:  fmt.Sprintf("字典编码格式错误: %s", codeStr),
			}, nil
		}
		dictCodes = append(dictCodes, code)
	}

	if len(dictCodes) == 0 {
		return &types.BaseResp{
			Code: 400,
			Msg:  "字典编码不能为空",
		}, nil
	}

	// 删除字典数据
	for _, dictCode := range dictCodes {
		err = l.svcCtx.SysDictDataModel.Delete(l.ctx, dictCode)
		if err != nil {
			l.Errorf("删除字典数据失败: dictCode=%d, err=%v", dictCode, err)
			return &types.BaseResp{
				Code: 500,
				Msg:  "删除字典数据失败",
			}, err
		}
	}

	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}, nil
}
