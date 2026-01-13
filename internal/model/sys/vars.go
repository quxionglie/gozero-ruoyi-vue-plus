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

// GetOffset 获取分页偏移量
// 返回 (PageNum - 1) * PageSize，确保结果 >= 0
func (p *PageQuery) GetOffset() int32 {
	if p == nil {
		return 0
	}
	offset := (p.PageNum - 1) * p.PageSize
	if offset < 0 {
		return 0
	}
	return offset
}

// GetOffsetAndLimit 获取分页偏移量和限制数量
// 返回 (offset, limit)，offset = (PageNum - 1) * PageSize，limit = PageSize
// offset 确保 >= 0
func (p *PageQuery) GetOffsetAndLimit() (int32, int32) {
	if p == nil {
		return 0, 10
	}
	offset := (p.PageNum - 1) * p.PageSize
	if offset < 0 {
		offset = 0
	}
	return offset, p.PageSize
}

// GetOrderDir 获取排序方向（兼容 asc、desc、ascending、descending 等）
// 返回 "asc" 或 "desc"，默认为 defaultDir（如果 defaultDir 为空则默认为 "asc"）
func (p *PageQuery) GetOrderDir(defaultDir string) string {
	if p == nil {
		if defaultDir == "" {
			return "asc"
		}
		return defaultDir
	}
	isAscStr := strings.ToLower(strings.TrimSpace(p.IsAsc))
	if isAscStr == "asc" || isAscStr == "ascending" {
		return "asc"
	} else if isAscStr == "desc" || isAscStr == "descending" {
		return "desc"
	}
	// 如果没有指定，使用默认值
	if defaultDir == "" {
		return "asc"
	}
	return defaultDir
}

// GetOrderBy 获取排序字段名（支持 camelCase 和 snake_case 转换，防止 SQL 注入）
// defaultColumn 是默认排序字段（如 "notice_id"）
// allowedOrderColumns 是允许的排序列映射（包含 snake_case 和 camelCase 两种格式）
// 返回处理后的排序字段名，如果输入的字段不在允许列表中，返回默认字段
func (p *PageQuery) GetOrderBy(defaultColumn string, allowedOrderColumns map[string]bool) string {
	if p == nil || p.OrderByColumn == "" {
		return defaultColumn
	}
	// 将 camelCase 转换为 snake_case
	columnName := camelToSnake(strings.TrimSpace(p.OrderByColumn))
	// 检查原始字段名和转换后的字段名是否在允许列表中
	originalColumn := strings.TrimSpace(p.OrderByColumn)
	if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
		// 使用转换后的 snake_case 字段名
		return columnName
	}
	// 如果不在允许列表中，返回默认字段
	return defaultColumn
}

// GetOrderByWithDir 获取排序字段名和方向（支持 camelCase 和 snake_case 转换，防止 SQL 注入）
// defaultOrderBy 是默认排序（如 "role_sort ASC, create_time ASC" 或 "oss_id ASC"）
// allowedOrderColumns 是允许的排序列映射（包含 snake_case 和 camelCase 两种格式）
// defaultDir 是默认排序方向（如 "asc"）
// 返回处理后的排序字符串（包含字段名和方向），如果输入的字段不在允许列表中，返回默认排序
func (p *PageQuery) GetOrderByWithDir(defaultOrderBy string, allowedOrderColumns map[string]bool, defaultDir string) string {
	if p == nil || p.OrderByColumn == "" {
		return defaultOrderBy
	}
	// 将 camelCase 转换为 snake_case
	columnName := camelToSnake(strings.TrimSpace(p.OrderByColumn))
	// 检查原始字段名和转换后的字段名是否在允许列表中
	originalColumn := strings.TrimSpace(p.OrderByColumn)
	if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
		// 使用转换后的 snake_case 字段名，并添加排序方向
		orderDir := p.GetOrderDir(defaultDir)
		return columnName + " " + strings.ToUpper(orderDir)
	}
	// 如果不在允许列表中，返回默认排序
	return defaultOrderBy
}

// GetOrderByWithDirAndPrefix 获取排序字段名和方向（支持表别名前缀、camelCase 和 snake_case 转换，防止 SQL 注入）
// defaultOrderBy 是默认排序（如 "u.user_id ASC"）
// allowedOrderColumns 是允许的排序列映射（包含表别名前缀，如 "u.user_id"）
// prefix 是表别名前缀（如 "u."），如果输入的字段没有前缀，会自动添加
// defaultDir 是默认排序方向（如 "asc"）
// 返回处理后的排序字符串（包含字段名和方向），如果输入的字段不在允许列表中，返回默认排序
func (p *PageQuery) GetOrderByWithDirAndPrefix(defaultOrderBy string, allowedOrderColumns map[string]bool, prefix string, defaultDir string) string {
	if p == nil || p.OrderByColumn == "" {
		return defaultOrderBy
	}
	originalColumn := strings.TrimSpace(p.OrderByColumn)
	// 如果没有表别名，添加前缀
	if !strings.Contains(originalColumn, ".") {
		originalColumn = prefix + originalColumn
	}
	// 将 camelCase 转换为 snake_case
	columnName := camelToSnake(originalColumn)
	// 检查原始字段名和转换后的字段名是否在允许列表中
	if allowedOrderColumns[originalColumn] || allowedOrderColumns[columnName] {
		// 使用转换后的 snake_case 字段名，并添加排序方向
		orderDir := p.GetOrderDir(defaultDir)
		return columnName + " " + strings.ToUpper(orderDir)
	}
	// 如果不在允许列表中，返回默认排序
	return defaultOrderBy
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

// snakeToCamel 将蛇形命名转换为驼峰命名
// 例如：login_time -> loginTime, user_name -> userName
func snakeToCamel(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, "_")
	var result strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result.WriteString(part)
		} else {
			// 首字母大写
			if len(part) > 0 {
				if part[0] >= 'a' && part[0] <= 'z' {
					result.WriteByte(part[0] - 32) // 转换为大写
				} else {
					result.WriteByte(part[0])
				}
				if len(part) > 1 {
					result.WriteString(part[1:])
				}
			}
		}
	}
	return result.String()
}

// buildAllowedOrderColumns 从字段名数组生成允许的排序列映射
// fieldNames 是带反引号的字段名数组，例如 ["`config_id`", "`tenant_id`"]
// 返回包含 snake_case 和 camelCase 两种格式的 map
func buildAllowedOrderColumns(fieldNames []string) map[string]bool {
	allowed := make(map[string]bool)
	for _, fieldName := range fieldNames {
		// 去除反引号
		snakeCase := strings.Trim(fieldName, "`")
		if snakeCase == "" {
			continue
		}
		// 添加 snake_case 格式
		allowed[snakeCase] = true
		// 添加 camelCase 格式
		camelCase := snakeToCamel(snakeCase)
		if camelCase != snakeCase {
			allowed[camelCase] = true
		}
	}
	return allowed
}

// buildAllowedOrderColumnsWithPrefix 从字段名数组生成带前缀的允许排序列映射
// fieldNames 是带反引号的字段名数组，例如 ["`config_id`", "`tenant_id`"]
// prefix 是表别名前缀，例如 "u."
// 返回包含 snake_case 和 camelCase 两种格式的 map
func buildAllowedOrderColumnsWithPrefix(fieldNames []string, prefix string) map[string]bool {
	allowed := make(map[string]bool)
	for _, fieldName := range fieldNames {
		// 去除反引号
		snakeCase := strings.Trim(fieldName, "`")
		if snakeCase == "" {
			continue
		}
		// 添加带前缀的 snake_case 格式
		allowed[prefix+snakeCase] = true
		// 添加带前缀的 camelCase 格式
		camelCase := snakeToCamel(snakeCase)
		if camelCase != snakeCase {
			allowed[prefix+camelCase] = true
		}
	}
	return allowed
}
