package gutils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ConfUrl struct {
	conf *url.URL
}

func NewConfUrl(conf string) ConfUrl {
	rconf, _ := url.Parse(conf)
	return ConfUrl{conf: rconf}
}
func (c ConfUrl) Scheme() string {
	return c.conf.Scheme
}
func (c ConfUrl) Host() string {
	return c.conf.Host
}
func (c ConfUrl) User() string {
	return c.conf.User.Username()
}
func (c ConfUrl) Pass() string {
	ps, _ := c.conf.User.Password()
	return ps
}
func (c ConfUrl) Path() string {
	r := c.conf.Path
	if strings.HasPrefix(r, "/") {
		return r[1:]
	}
	return r
}
func (c ConfUrl) Port() int {
	p := c.conf.Port()
	rs, e := strconv.Atoi(p)
	if e != nil {
		return 0
	}
	return rs

}
func (c ConfUrl) Get(key, defv string) string {
	if !c.conf.Query().Has(key) {
		return defv
	}
	return c.conf.Query().Get(key)
}
func (c ConfUrl) GetInt(key string, defv int) int {
	r := c.Get(key, fmt.Sprintf("%d", defv))
	rs, e := strconv.Atoi(r)
	if e != nil {
		return defv
	}
	return rs
}
func (c ConfUrl) GetBool(key string, defv bool) bool {
	r := c.Get(key, fmt.Sprintf("%v", defv))
	rs, e := strconv.ParseBool(r)
	if e != nil {
		return defv
	}
	return rs
}
func (c ConfUrl) GetDuration(key string, defv int) time.Duration {
	r := c.GetInt(key, defv)
	return time.Duration(int64(r)) * time.Second
}
