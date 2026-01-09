package util

import "gozero-ruoyi-vue-plus/internal/types"

// Ok 成功响应（返回 BaseResp，兼容新类型）
func Ok(data interface{}) *types.BaseResp {
	return &types.BaseResp{
		Code: 200,
		Msg:  "操作成功",
	}
}

// Fail 失败响应
func Fail(msg string) *types.BaseResp {
	return &types.BaseResp{
		Code: 500,
		Msg:  msg,
	}
}

// FailWithCode 带错误码的失败响应
func FailWithCode(code int32, msg string) *types.BaseResp {
	return &types.BaseResp{
		Code: code,
		Msg:  msg,
	}
}
