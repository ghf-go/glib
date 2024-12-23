package alipay

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

// 正式环境：https://openapi.alipay.com
// 沙箱环境：https://openapi-sandbox.dl.alipaydev.com

type Client struct{}

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
