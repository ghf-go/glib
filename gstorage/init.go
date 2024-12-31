package gstorage

import "time"

type Storage interface {
	//获取上传token
	GetToken(fileKey, fileName string) map[string]any
}

func getGMTISO8601(expireEnd int64) string {
	return time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
}
