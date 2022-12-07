package Rrpc

import "strings"

func SubStringLast(str string, substr string) string {
	//先查找有没有
	index := strings.Index(str, substr)
	if index == -1 {
		return ""
	}
	len := len(substr)
	return str[index+len:]
}
