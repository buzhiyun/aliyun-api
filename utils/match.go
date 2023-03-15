package utils

import "strings"

// matchWildcard 判断字符串 s 是否符合通配符模式 pattern
func MatchWildcard(s string, pattern string) bool {
	if pattern == "" {
		return s == ""
	}
	// 将 pattern 按 * 分割成多个子模式
	patterns := strings.Split(pattern, "*")
	// 如果只有一个子模式，直接使用 strings.HasPrefix 和 strings.HasSuffix 判断
	if len(patterns) == 1 {
		return strings.HasPrefix(s, patterns[0]) && strings.HasSuffix(s, patterns[0])
	}
	// 处理第一个子模式前面的部分
	if !strings.HasPrefix(s, patterns[0]) {
		return false
	}
	// 处理最后一个子模式后面的部分
	if !strings.HasSuffix(s, patterns[len(patterns)-1]) {
		return false
	}
	// 处理中间的子模式
	lastIndex := len(patterns[0])
	for i := 0; i < len(patterns)-1; i++ {
		p := patterns[i+1]
		if p == "" {
			continue
		}
		index := strings.Index(s[lastIndex:], p)
		if index == -1 {
			return false
		}
		lastIndex += index + len(p)
	}
	return true
}
