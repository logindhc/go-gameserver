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

func StringToMap(input string) map[int]int {
	// 去掉花括号
	input = input[1 : len(input)-1]

	// 初始化结果 map
	result := make(map[int]int)

	// 遍历字符串
	i := 0
	for i < len(input) {
		// 找到键的开始位置
		startKey := i
		for i < len(input) && unicode.IsDigit(rune(input[i])) {
			i++
		}
		if i == startKey {
			break // 没有找到有效的键
		}

		// 解析键
		keyStr := input[startKey:i]
		key, err := strconv.Atoi(keyStr)
		if err != nil {
			break // 键解析失败
		}

		// 跳过冒号
		if i < len(input) && input[i] == ':' {
			i++
		} else {
			break // 没有找到冒号
		}

		// 找到值的结束位置
		startValue := i
		for i < len(input) && unicode.IsDigit(rune(input[i])) {
			i++
		}
		if i == startValue {
			break // 没有找到有效的值
		}

		// 解析值
		valueStr := input[startValue:i]
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			break // 值解析失败
		}

		// 存储键值对
		result[key] = value

		// 跳过逗号
		if i < len(input) && input[i] == ',' {
			i++
		}
	}

	return result
}
