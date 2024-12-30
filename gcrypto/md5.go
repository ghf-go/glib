package gcrypto

import (
	"crypto/md5"
	"fmt"
)

func MD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
func MD5Bytes(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
