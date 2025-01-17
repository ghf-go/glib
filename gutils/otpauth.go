package gutils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strconv"
	"time"
)

type otp2Fa struct {
	secret string
}

// otpauth://totp/GitHub:eoe2005?secret=16bit&issuer=GitHub
// 验证
func VerifyOtp2Fa(secret, code string) bool {

	if ci, e := strconv.ParseInt(code, 10, 32); e == nil {
		o := &otp2Fa{
			secret: secret,
		}
		return o.VerifyCode(int32(ci))
	}
	return false
}

// 为了考虑时间误差，判断前当前时间及前后30秒时间
func (o *otp2Fa) VerifyCode(code int32) bool {
	// 当前google值
	if o.getCode(0) == code {
		return true
	}

	// 前30秒google值
	if o.getCode(-30) == code {
		return true
	}

	// 后30秒google值
	if o.getCode(30) == code {
		return true
	}

	return false
}

// 获取Google Code
func (o *otp2Fa) getCode(offset int64) int32 {
	key, err := base32.StdEncoding.DecodeString(o.secret)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix() + offset
	return int32(o.oneTimePassword(key, o.toBytes(epochSeconds/30)))
}

// from https://github.com/robbiev/two-factor-auth/blob/master/main.go
func (o *otp2Fa) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func (o *otp2Fa) toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func (o *otp2Fa) oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := o.toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}
