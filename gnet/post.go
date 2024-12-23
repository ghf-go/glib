package gnet

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// http Post 请求
func PostContent(url, contentType, body string) ([]byte, error) {
	r, e := http.Post(url, contentType, strings.NewReader(body))
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// http Post 请求
func PostFormContent(qurl string, data url.Values) ([]byte, error) {
	r, e := http.PostForm(qurl, data)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

// http Post 获取JSON
func PostJson(url, contentType, body string, out any) error {
	data, e := PostContent(url, contentType, body)
	if e != nil {
		return e
	}
	return json.Unmarshal(data, out)
}

// http PostForm 请求
func PostFormJson(qurl string, rdata url.Values, out any) error {
	data, e := PostFormContent(qurl, rdata)
	if e != nil {
		return e
	}
	return json.Unmarshal(data, out)
}
