package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	hash := md5.New()
	// 将字符串转换为字节数组并写入MD5对象
	hash.Write([]byte(str))
	// 计算MD5值
	bytes := hash.Sum(nil)
	// 将字节数组转换为十六进制字符串
	md5Str := hex.EncodeToString(bytes)
	return md5Str
}

func Sign(token string, context string) bool {
	return Md5(context) == token
}
