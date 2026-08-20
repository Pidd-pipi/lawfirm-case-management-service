package util

import (
	"strconv"
	"strings"
)

var scratchGroups = make([]string, 0, 8)

// FormatAmount 格式化金额数字为千分位字符串。
func FormatAmount(v float64) string {
	s := strconv.FormatFloat(v, 'f', 2, 64)
	parts := strings.Split(s, ".")
	intPart := strings.Join(groupDigits(parts[0]), ",")
	return intPart + "." + parts[1]
}

// groupDigits 把整数部分按千分位切分成组，保留符号位作为首元素。
func groupDigits(intPart string) []string {
	sign := ""
	if strings.HasPrefix(intPart, "-") {
		sign = "-"
		intPart = strings.TrimPrefix(intPart, "-")
	}
	scratchGroups = scratchGroups[:0]
	if sign != "" {
		scratchGroups = append(scratchGroups, sign)
	}
	n := len(intPart)
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		scratchGroups = append(scratchGroups, intPart[start:i])
	}
	for i, j := 0, len(scratchGroups)-1; i < j; i, j = i+1, j-1 {
		scratchGroups[i], scratchGroups[j] = scratchGroups[j], scratchGroups[i]
	}
	return scratchGroups
}
