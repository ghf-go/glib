package gstorage

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

type qiniuStorage struct {
	client     *auth.Credentials
	bucket     string
	cdnHost    string
	uploadHost string
	ak         string
	sk         string
}

func NewQiniu(conf string) *qiniuStorage {
	rconf, e := url.Parse(conf)
	if e != nil {
		panic(e.Error())
	}
	upass, _ := rconf.User.Password()
	return &qiniuStorage{
		client:     credentials.NewCredentials(rconf.User.Username(), upass),
		bucket:     rconf.Query().Get("bucket"),
		cdnHost:    rconf.Query().Get("cdn_host"),
		uploadHost: rconf.Query().Get("upload_host"),
		ak:         rconf.User.Username(),
		sk:         upass,
	}
}
func (s *qiniuStorage) GetToken(fileKey, fileName string) map[string]any {
	ext := path.Ext(fileName)
	fileName = path.Base(fileName)
	fkey := fmt.Sprintf("%s/%s%s", time.Now().Format("2006/01/02"), fileKey, ext)
	putPolicy, err := uptoken.NewPutPolicyWithKey(s.bucket, fkey, time.Now().Add(1*time.Hour))
	if err != nil {
		return map[string]any{}
	}

	putPolicy.SetReturnBody(`{"key":"$(key)","url":"$(x:url)","file_size":$(fsize),"file_key":"$(x:file_key)","file_name":"$(x:file_name)"}`)
	upToken, err := uptoken.NewSigner(putPolicy, s.client).GetUpToken(context.Background())
	if err != nil {
		return map[string]any{}
	}

	return map[string]any{
		"driver":      "qiniu",
		"upload_host": s.uploadHost,
		"data": map[string]any{
			"token":       upToken,
			"key":         fkey,
			"x:file_key":  fileKey,
			"x:file_name": fileName,
			"x:url":       s.cdnHost + fkey,
		},
	}
}
