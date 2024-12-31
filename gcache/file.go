package gcache

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ghf-go/glib/gcrypto"
	"github.com/ghf-go/glib/gutils"
)

func NewFileCache(c gutils.ConfUrl) *fileCache {
	return &fileCache{
		dirName: c.Get("dir", "/tmp"),
	}
}

type fileCache struct {
	dirName string
}
type fileCacheContent struct {
	ExpireAt int64  `json:"expire"`
	Content  string `json:"data"`
}

func (c *fileCache) Get(key string, defval string) string {
	file := c.getFileName(key)
	data, e := os.ReadFile(file)
	if e != nil {
		return defval
	}
	val := &fileCacheContent{}
	e = json.Unmarshal(data, val)
	if e != nil {
		return defval
	}
	if val.ExpireAt == 0 || val.ExpireAt > time.Now().Unix() {
		return val.Content
	}
	return defval
}
func (c *fileCache) GetAll(key ...string) map[string]string {
	ret := map[string]string{}
	for _, k := range key {
		ret[k] = c.Get(k, "")
	}
	return ret
}
func (c *fileCache) GetObj(key string, out any) error {
	return json.Unmarshal([]byte(c.Get(key, "")), out)
}
func (c *fileCache) GetAllObj(data map[string]any) {
	for k, v := range data {
		e := c.GetObj(k, v)
		if e == nil {
			data[k] = v
		}
	}
}
func (c *fileCache) Set(key, val string, timeOut ...int) error {
	tt := 0
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	fn := c.getFileName(key)
	dir := path.Dir(fn)
	os.MkdirAll(dir, os.ModePerm)
	v := &fileCacheContent{
		Content:  val,
		ExpireAt: 0,
	}
	if tt > 0 {
		v.ExpireAt = time.Now().Unix() + int64(tt)
	}
	wd, e := json.Marshal(v)
	if e != nil {
		return e
	}
	os.WriteFile(fn, wd, os.ModePerm)
	return nil
}
func (c *fileCache) SetObj(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.Set(key, string(wd), timeOut...)
}
func (c *fileCache) SetNx(key, val string, timeOut ...int) error {
	dv := "-not-exits-"
	vv := c.Get(key, dv)
	if dv == vv {
		return fmt.Errorf("%s 已存在", key)
	}
	return c.Set(key, val, timeOut...)
}
func (c *fileCache) SetObjNx(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.SetNx(key, string(wd), timeOut...)
}
func (c *fileCache) Del(key ...string) error {
	for _, k := range key {
		fn := c.getFileName(k)
		os.Remove(fn)
	}
	return nil
}
func (c *fileCache) Flush() error {
	ds, e := os.ReadDir(c.dirName)
	if e != nil {
		return e
	}
	for _, v := range ds {
		os.RemoveAll(c.dirName + "/" + v.Name())
	}
	return nil
}
func (c *fileCache) Lock(key string, callfunc func(), timeOut ...int) error {
	tt := 30
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	e := c.SetNx(key, "lock", tt)
	if e != nil {
		return e
	}
	callfunc()
	return c.Del(key)
}
func (c *fileCache) Incr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	ret := c.Get(key, "0")
	rv, e := strconv.Atoi(ret)
	if e != nil {
		rv = sv
	} else {
		rv += sv
	}
	c.Set(key, fmt.Sprintf("%d", rv))
	return rv
}
func (c *fileCache) Decr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	ret := c.Get(key, "0")
	rv, e := strconv.Atoi(ret)
	if e != nil {
		rv = -sv
	} else {
		rv -= sv
	}
	c.Set(key, fmt.Sprintf("%d", rv))
	return rv
}

// 获取文件路径
func (c *fileCache) getFileName(key string) string {
	md5str := gcrypto.MD5(key)
	return fmt.Sprintf("%s/%s/%s/%s/%s", c.dirName, md5str[0:2], md5str[2:4], md5str[4:6], md5str[6:])
}
