package utils

import "strings"

//IsEmpty 是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
