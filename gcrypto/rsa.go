package gcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// 从字符串中解析私钥
func RsaGetPrivatekey(keyData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return privKey, nil
}

// 从字符串中解析公钥钥
func RsaGetPublicKey(keyData string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	pk, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return pk, nil
}

// 使用给定私钥对数据进行 rsa 签名
func RsaSignString(pkData, data string) (string, error) {
	pk, e := RsaGetPrivatekey(pkData)
	if e != nil {
		return "", fmt.Errorf("failed to parse private key: %v", e)
	}
	r, e := RsaSign(pk, []byte(data))
	if e != nil {
		return "", fmt.Errorf("failed to parse private key: %v", e)
	}
	return string(r), nil
}

// 使用给定私钥对数据进行 rsa 签名
func RsaSign(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to hash: %v", err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}

	return signature, nil
}

// 验证给定公钥和签名的数据
func RsaVerifyString(pkData, data, signature string) error {
	pk, e := RsaGetPublicKey(pkData)
	if e != nil {
		return fmt.Errorf("failed to parse public key: %v", e)
	}
	return RsaVerify(pk, []byte(data), []byte(signature))

}

// 验证给定公钥和签名的数据
func RsaVerify(pubKey *rsa.PublicKey, data, signature []byte) error {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to hash: %v", err)
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash.Sum(nil), signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	return nil
}
