package alipay

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/ghf-go/glib/gcrypto"
)

// 正式环境：https://openapi.alipay.com
// 沙箱环境：https://openapi-sandbox.dl.alipaydev.com

type Client struct {
	isDebug          bool
	appID            string //
	appCertSn        string
	alipayRootCertSn string
	appPrivateKey    *rsa.PrivateKey
	notifyUrl        string
}

// NewDebugClient 创建一个用于调试的支付宝客户端实例
// 参数:
//   - appID: 支付宝应用ID
//   - notifyUrl: 支付宝异步通知URL
//   - appKey: 应用私钥
//   - rootkey: 支付宝根证书内容
//
// 返回值:
//   - *Client: 支付宝客户端实例
//   - error: 错误信息，如果有的话
func NewDebugClient(appID, notifyUrl, appKey, rootkey string) (*Client, error) {
	ret := &Client{
		isDebug:   true,
		notifyUrl: notifyUrl,
		appID:     appID,
	}
	p, e := gcrypto.RsaGetPrivatekey(appKey)
	if e != nil {
		return nil, e
	}
	ret.appPrivateKey = p
	if rootkey != "" {
		ret.appCertSn = getCertSN(appKey)
		ret.alipayRootCertSn = getRootCertSN(rootkey)
	}
	return ret, nil
}
func NewClient(appID, notifyUrl, appKey, rootkey string) (*Client, error) {
	ret := &Client{
		isDebug:   false,
		notifyUrl: notifyUrl,
		appID:     appID,
	}
	p, e := gcrypto.RsaGetPrivatekey(appKey)
	if e != nil {
		return nil, e
	}
	ret.appPrivateKey = p
	if rootkey != "" {
		ret.appCertSn = getCertSN(appKey)
		ret.alipayRootCertSn = getRootCertSN(rootkey)
	}
	return ret, nil
}

func (c *Client) getApiHost() string {
	if c.isDebug {
		return "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	}
	return "https://openapi.alipay.com/gateway.do"
}

// 获取应用sn
func getCertSN(certContent string) string {
	p, _ := pem.Decode([]byte(certContent))
	c, e := x509.ParseCertificate(p.Bytes)
	if e != nil {
		fmt.Println(e.Error())
		panic(e.Error())
	}
	dd := md5.Sum([]byte(fmt.Sprintf("%v%d", c.Issuer, c.SerialNumber)))
	return fmt.Sprintf("%x", dd)
}

// 获取根证书信息
func getRootCertSN(rootCertContent string) string {
	arrays := strings.Split(rootCertContent, "-----END CERTIFICATE-----")
	ret := ""
	for _, item := range arrays {
		if ret == "" {
			ret = getCertSN(item + "-----END CERTIFICATE-----")
		} else {
			ret += "_" + getCertSN(item+"-----END CERTIFICATE-----")
		}
	}
	return ret
}
