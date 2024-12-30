package gutils

import "regexp"

// 是否是手机号
func IsMobile(name string) bool {
	return regexp.MustCompile(`^1\d{9,}$`).MatchString(name)
}

// 是否是邮箱
func IsEmail(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(name)
}
