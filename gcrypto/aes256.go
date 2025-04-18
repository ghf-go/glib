package gcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 加密函数
func Aes256Encrypt(plaintext []byte, key []byte) (string, error) {
	// 检查密钥长度是否为 32 字节（AES-256）
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}

	// 生成随机的 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 创建 AES 加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 对明文进行 PKCS#7 填充
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// 创建 CBC 加密器
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// 将 IV 和密文拼接在一起
	result := append(iv, ciphertext...)

	// 返回 Base64 编码后的结果
	return base64.StdEncoding.EncodeToString(result), nil
}

// 解密函数
func Aes256Decrypt(base64Ciphertext string, key []byte) (string, error) {
	// 检查密钥长度是否为 32 字节（AES-256）
	if len(key) != 32 {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}

	// Base64 解码
	ciphertext, err := base64.StdEncoding.DecodeString(base64Ciphertext)
	if err != nil {
		return "", err
	}

	// 检查密文长度是否有效
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// 提取 IV 和实际密文
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建 AES 解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 CBC 解密器
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 移除 PKCS#7 填充
	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", err
	}

	// 返回解密后的字符串
	return string(plaintext), nil
}

// PKCS#7 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// 移除 PKCS#7 填充
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 || length%blockSize != 0 {
		return nil, errors.New("invalid data length")
	}
	padding := int(data[length-1])
	if padding > blockSize || padding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}

// 主函数
// func main() {
// 	// 密钥（必须是 32 字节）
// 	key := []byte("this-is-a-32-byte-key-for-aes-256-cbc")

// 	// 明文
// 	plaintext := "Hello, AES-256-CBC!"

// 	// 加密
// 	encrypted, err := encrypt([]byte(plaintext), key)
// 	if err != nil {
// 		fmt.Println("Encryption error:", err)
// 		return
// 	}
// 	fmt.Println("Encrypted:", encrypted)

// 	// 解密
// 	decrypted, err := decrypt(encrypted, key)
// 	if err != nil {
// 		fmt.Println("Decryption error:", err)
// 		return
// 	}
// 	fmt.Println("Decrypted:", decrypted)
// }
