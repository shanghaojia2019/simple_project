package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

//MD5  生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

//Md5 MD5加密
func Md5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
