package gcache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ghf-go/glib/gutils"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	rd   *redis.Client
	conf string
	ctx  context.Context
}

func NewRedisClient(conf gutils.ConfUrl) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     conf.Host(),
		Username: rconf.User.Username(),
		Password: upass,
		// MinIdleConns:    rconf.Query().Get(),
		// MaxIdleConns:    rconf.MaxIdleConns,
		// MaxActiveConns:  rconf.MaxActiveConns,
		// ConnMaxIdleTime: time.Minute * time.Duration(rconf.ConnMaxIdleTime),
		// ConnMaxLifetime: time.Minute * time.Duration(rconf.ConnMaxLifetime),
	})
}

func NewRedisCache(conf string) *redisCache {
	return &redisCache{
		rd:   NewRedisClient(conf),
		conf: conf,
		ctx:  context.Background(),
	}
}
func (c *redisCache) Get(key string, defval string) string {
	ret, e := c.rd.Get(c.ctx, key).Result()
	if e != nil {
		return defval
	}
	return ret
}
func (c *redisCache) GetAll(key ...string) map[string]string {
	rs := c.rd.MGet(c.ctx, key...).Val()

	ret := map[string]string{}
	for i, v := range rs {
		ret[key[i]] = v.(string)
	}
	return ret
}
func (c *redisCache) GetObj(key string, out any) error {
	return json.Unmarshal([]byte(c.Get(key, "")), out)
}
func (c *redisCache) GetAllObj(data map[string]any) {
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
func (c *redisCache) Set(key, val string, timeOut ...int) error {
	tt := 0
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	return c.rd.Set(c.ctx, key, val, time.Second*time.Duration(int64(tt))).Err()
}
func (c *redisCache) SetObj(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.Set(key, string(wd), timeOut...)
}
func (c *redisCache) SetNx(key, val string, timeOut ...int) error {
	tt := 30
	if len(timeOut) > 0 {
		tt = timeOut[0]
	}
	r, e := c.rd.SetNX(c.ctx, key, "", time.Duration(int64(tt))*time.Second).Result()
	if e != nil {
		return e
	}
	if !r {
		return fmt.Errorf("%s 存在", key)
	}
	return nil

}
func (c *redisCache) SetObjNx(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.SetNx(key, string(wd), timeOut...)
}
func (c *redisCache) Incr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.rd.IncrBy(c.ctx, key, int64(sv)).Result()
	if e != nil {
		return sv
	}
	return int(r)
}
func (c *redisCache) Decr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.rd.DecrBy(c.ctx, key, int64(sv)).Result()
	if e != nil {
		return -sv
	}
	return int(r)
}
func (c *redisCache) Del(key ...string) error {
	return c.rd.Del(c.ctx, key...).Err()
}
func (c *redisCache) Flush() error {
	return c.rd.FlushDB(c.ctx).Err()
}
func (c *redisCache) Lock(key string, callfunc func(), timeOut ...int) error {
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
