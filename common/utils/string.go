package utils

import (
	"strconv"
	"unicode"
)

func StrToInt64(str string) (num int64) {
	num, _ = strconv.ParseInt(str, 10, 64)
	return
}
func StrToInt(str string) (num int) {
	num, _ = strconv.Atoi(str)
	return
}

// CamelToSnake 将驼峰命名转换为下划线命名
func CamelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && !(s == "ID" && i == 1) { // 特殊处理 ID
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
