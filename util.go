package Rrpc

import (
	"strings"
	"unicode"
	"unsafe"
)

func SubStringLast(str string, substr string) string {
	//先查找有没有
	index := strings.Index(str, substr)
	if index == -1 {
		return ""
	}
	len := len(substr)
	return str[index+len:]
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
