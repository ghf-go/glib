package gutils

import (
	"fmt"
	"strings"

	"golang.org/x/exp/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 生成随机字符串
func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

// 隐藏手机号
func HideMobile(mobile string) string {
	return fmt.Sprintf("%s****%s", mobile[0:3], mobile[7:])
}

// 隐藏邮箱
func HideEmail(email string) string {
	i := strings.Index(email, "@")
	name := email[:i]
	host := email[i:]
	ln := len(name)
	if ln == 1 {
		return fmt.Sprintf("%s****%s", name, host)
	} else if ln > 5 {
		return fmt.Sprintf("%s****%s%s", name[:3], name[ln-3:], host)
	} else {
		return fmt.Sprintf("%s****%s", name[:ln-2], host)
	}

}
