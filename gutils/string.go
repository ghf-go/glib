package gutils

import (
	"fmt"
	"strconv"
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

func String2Uint64(src string) uint64 {
	r, e := strconv.ParseUint(src, 10, 64)
	if e != nil {
		return 0
	}
	return r

}
func String2Uint32(src string) uint32 {
	r, e := strconv.ParseUint(src, 10, 32)
	if e != nil {
		return 0
	}
	return uint32(r)
}

func String2Uint16(src string) uint16 {
	r, e := strconv.ParseUint(src, 10, 16)
	if e != nil {
		return 0
	}
	return uint16(r)
}
func String2Uint(src string) uint {
	r, e := strconv.ParseUint(src, 10, 64)
	if e != nil {
		return 0
	}
	return uint(r)
}

func String2Uint8(src string) uint8 {
	r, e := strconv.ParseUint(src, 10, 8)
	if e != nil {
		return 0
	}
	return uint8(r)
}
