package gcache

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/ghf-go/glib/gutils"
)

type memcacheCache struct {
	client *memcache.Client
	conf   string
}

func NewMemcache(conf gutils.ConfUrl) *memcacheCache {
	return &memcacheCache{
		client: memcache.New(strings.Split(conf, ",")...),
		conf:   conf,
	}
}
func (c *memcacheCache) Get(key string, defval string) string {
	val, e := c.client.Get(key)
	if e != nil {
		return defval
	}
	return string(val.Value)
}
func (c *memcacheCache) GetAll(key ...string) map[string]string {
	rs, _ := c.client.GetMulti(key)
	ret := map[string]string{}
	for _, k := range key {
		v, ok := rs[k]
		if ok {
			ret[k] = string(v.Value)
		} else {
			ret[k] = ""
		}
	}
	return ret
}
func (c *memcacheCache) GetObj(key string, out any) error {
	return json.Unmarshal([]byte(c.Get(key, "")), out)
}
func (c *memcacheCache) GetAllObj(data map[string]any) {
	keys := []string{}
	for k, _ := range data {
		keys = append(keys, k)
	}
	retStr := c.GetAll(keys...)
	for k, v := range retStr {
		if vv, ok := data[k]; ok {
			json.Unmarshal([]byte(v), vv)
			data[k] = vv
		}
	}
}
func (c *memcacheCache) Set(key, val string, timeOut ...int) error {
	tt := 0
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(val),
		Expiration: int32(tt),
	}

	return c.client.Set(item)
}
func (c *memcacheCache) SetObj(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.Set(key, string(wd), timeOut...)
}
func (c *memcacheCache) SetNx(key, val string, timeOut ...int) error {
	rv := "-not-exitxs-key-"
	if c.Get(key, rv) == rv {
		return fmt.Errorf("%s 存在", key)
	}
	return c.Set(key, val, timeOut...)
}
func (c *memcacheCache) SetObjNx(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.SetNx(key, string(wd), timeOut...)
}
func (c *memcacheCache) Incr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.client.Increment(key, uint64(sv))
	if e != nil {
		return sv
	}
	return int(r)
}
func (c *memcacheCache) Decr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.client.Decrement(key, uint64(sv))
	if e != nil {
		return -sv
	}
	return int(r)
}
func (c *memcacheCache) Del(key ...string) error {
	for _, k := range key {
		c.client.Delete(k)
	}
	return nil
}
func (c *memcacheCache) Flush() error {
	return c.client.FlushAll()
}
func (c *memcacheCache) Lock(key string, callfunc func(), timeOut ...int) error {
	tt := 30
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	e := c.SetNx(key, "", tt)
	if e != nil {
		return e
	}
	callfunc()
	return c.Del(key)
}
