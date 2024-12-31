package gstorage

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path"
	"time"
)

// https://github.com/aliyun/alibabacloud-oss-go-sdk-v2/blob/master/README-CN.md
// https://help.aliyun.com/zh/oss/user-guide/simple-upload?spm=a2c4g.11186623.help-menu-31815.d_2_3_1_0.658930244Xpq3w
// https://help.aliyun.com/zh/oss/user-guide/form-upload?spm=a2c4g.11186623.help-menu-31815.d_2_3_1_4.5d4cb415B70GWk#22bb72dc11p2g
type aliossStorage struct {
	bucket     string
	cdnHost    string
	uploadHost string
	ak         string
	sk         string
}

func NewAliOss(conf string) *aliossStorage {
	rconf, e := url.Parse(conf)
	if e != nil {
		panic(e.Error())
	}
	upass, _ := rconf.User.Password()
	return &aliossStorage{
		bucket:     rconf.Query().Get("bucket"),
		cdnHost:    rconf.Query().Get("cdn_host"),
		uploadHost: rconf.Query().Get("upload_host"),
		ak:         rconf.User.Username(),
		sk:         upass,
	}
}

type aliOssConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

func (s *aliossStorage) GetToken(fileKey, fileName string) map[string]any {
	ext := path.Ext(fileName)
	fileName = path.Base(fileName)
	uploadDir := time.Now().Format("2006/01/02") + "/"
	fkey := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), fileKey, ext)
	now := time.Now().Unix()
	expireEnd := now + int64(3600)
	tokenExpire := getGMTISO8601(expireEnd)
	var config aliOssConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)
	result, err := json.Marshal(config)
	if err != nil {
		fmt.Println("callback json err:", err)
		return map[string]any{}
	}
	encodedResult := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(sha1.New, []byte(s.sk))
	io.WriteString(h, encodedResult)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return map[string]any{
		"driver":         "alioss",
		"ossAccessKeyId": s.sk,
		"host":           fmt.Sprintf("http://%s.%s", s.bucket, s.uploadHost),
		"signature":      signedStr,
		"policy":         encodedResult,
		"dir":            uploadDir,
		"key":            fkey,
	}

}
