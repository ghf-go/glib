package gutils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type confUrl struct {
	conf *url.URL
}

func NewConfUrl(conf string) confUrl {
	rconf, _ := url.Parse(conf)
	return confUrl{conf: rconf}
}
func (c confUrl) Scheme() string {
	return c.conf.Scheme
}
func (c confUrl) Host() string {
	return c.conf.Host
}
func (c confUrl) User() string {
	return c.conf.User.Username()
}
func (c confUrl) Pass() string {
	ps, _ := c.conf.User.Password()
	return ps
}
func (c confUrl) Path() string {
	r := c.conf.Path
	if strings.HasPrefix(r, "/") {
		return r[1:]
	}
	return r
}
func (c confUrl) Get(key, defv string) string {
	if !c.conf.Query().Has(key) {
		return defv
	}
	return c.conf.Query().Get(key)
}
func (c confUrl) GetInt(key string, defv int) int {
	r := c.Get(key, fmt.Sprintf("%d", defv))
	rs, e := strconv.Atoi(r)
	if e != nil {
		return defv
	}
	return rs
}
func (c confUrl) GetBool(key string, defv bool) bool {
	r := c.Get(key, fmt.Sprintf("%v", defv))
	rs, e := strconv.ParseBool(r)
	if e != nil {
		return defv
	}
	return rs
}
func (c confUrl) GetDuration(key string, defv int) time.Duration {
	r := c.GetInt(key, defv)
	return time.Duration(int64(r)) * time.Second
}
