package gnet

import (
	"encoding/json"
	"io"
	"net/http"
)

// http Get 请求
func GetContent(url string) ([]byte, error) {
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// http Get 获取JSON
func GetJson(url string, out any) error {
	data, e := GetContent(url)
	if e != nil {
		return e
	}
	return json.Unmarshal(data, out)
}
