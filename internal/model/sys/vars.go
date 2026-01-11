package sys

import (
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// camelToSnake 将驼峰命名转换为蛇形命名（简化版）
// 例如：loginTime -> login_time, userName -> user_name
func camelToSnake(s string) string {
	if s == "" {
		return s
	}
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
			result.WriteRune(r + 32) // 转换为小写
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(r + 32) // 转换为小写
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
