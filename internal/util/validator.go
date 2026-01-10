package util

import (
	"fmt"
	"regexp"
)

// ValidateStringLength 校验字符串长度
func ValidateStringLength(value, fieldName string, maxLen int) error {
	if len([]rune(value)) > maxLen {
		return fmt.Errorf("%s长度不能超过%d个字符", fieldName, maxLen)
	}
	return nil
}

// ValidateDictType 校验字典类型格式：必须以小写字母开头，只能包含小写字母、数字、下划线
func ValidateDictType(dictType string) error {
	// 正则表达式：^[a-z][a-z0-9_]*$
	matched, err := regexp.MatchString("^[a-z][a-z0-9_]*$", dictType)
	if err != nil {
		return fmt.Errorf("字典类型格式校验失败: %v", err)
	}
	if !matched {
		return fmt.Errorf("字典类型必须以字母开头，且只能为（小写字母，数字，下划线）")
	}
	return nil
}
