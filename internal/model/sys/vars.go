package sys

import (
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// PageQuery 分页查询参数
type PageQuery struct {
	PageNum       int32  // 当前页数
	PageSize      int32  // 分页大小
	OrderByColumn string // 排序列
	IsAsc         string // 排序方向（desc 或 asc）
}

// Normalize 初始化分页参数的非合规值
// pageNum <= 0 时设为 1，pageSize <= 0 时设为 10
func (p *PageQuery) Normalize() {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

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
