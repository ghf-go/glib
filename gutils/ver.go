package gutils

import (
	"math"
	"strconv"
	"strings"
)

// ver1 的版本是否大于ver2
func CheckVersion(ver1, ver2 string) bool {
	vs1 := strings.Split(ver1, ".")
	vs2 := strings.Split(ver2, ".")
	ls := max(len(vs1), len(vs2))
	v1 := float64(0)
	v2 := float64(0)
	for i, v := range vs1 {
		vv, _ := strconv.ParseFloat(v, 10)
		v1 += math.Pow10((ls-i-1)*2) * vv
	}
	for i, v := range vs2 {
		vv, _ := strconv.ParseFloat(v, 10)
		v2 += math.Pow10((ls-i-1)*2) * vv
	}
	return v1 > v2
}
