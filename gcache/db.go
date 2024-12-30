package gcache

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type dbCache struct {
	db sql.DB
	tb string
}

func NewDbCache(db sql.DB) *dbCache {
	tb := "t_cache"
	rtb := ""
	db.QueryRow("show tables like ?", tb).Scan(&rtb)
	if tb != rtb {
		db.Exec(fmt.Sprintf("CREATE TABLE`%s` (`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,`key` VARCHAR(30) NOT NULL DEFAULT '' COMMENT 'KEY',`val` TEXT NOT NULL COMMENT 'VAL',`expire_at` INT (11) NOT NULL DEFAULT 0 COMMENT '到期时间戳', PRIMARY KEY (`id`),UNIQUE KEY `uniq_key` (`key`),) ENGINE = innodb DEFAULT CHARSET = utf8mb4 COMMENT = '缓存表';", tb))
	}
	return &dbCache{
		db: db,
		tb: tb,
	}
}

func (c *dbCache) Get(key string, defval string) string {
	val := ""
	expire := 0
	e := c.db.QueryRow(fmt.Sprintf("SELECT val,expire_at FROM %s WHERE key=?", c.tb), key).Scan(&val, &expire)
	if e != nil {
		return defval
	}
	if expire == 0 || expire > int(time.Now().Unix()) {
		return val
	}
	return defval
}
func (c *dbCache) GetAll(key ...string) map[string]string {
	ret := map[string]string{}
	for _, k := range key {
		ret[k] = ""
	}
	rows, e := c.db.Query(fmt.Sprintf("SELECT key,val,expire_at FROM %s WHERE key IN (?)", c.tb), strings.Join(key, "','"))
	if e != nil {
		return ret
	}
	for rows.Next() {
		key := ""
		val := ""
		expire := 0
		rows.Scan(&key, &val, &expire)
		if expire == 0 || expire > int(time.Now().Unix()) {
			ret[key] = val
		}
	}
	return ret
}
func (c *dbCache) GetObj(key string, out any) error {
	return json.Unmarshal([]byte(c.Get(key, "")), out)
}
func (c *dbCache) GetAllObj(data map[string]any) {
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
func (c *dbCache) Set(key, val string, timeOut ...int) error {
	tt := 0
	if len(timeOut) > 0 {
		tt = timeOut[0] + int(time.Now().Unix())
	}
	_, e := c.db.Exec(fmt.Sprintf("REPLACE %s(key,val,expire_at) VALUES(?,?,?)", c.tb), key, val, tt)
	return e
}
func (c *dbCache) SetObj(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.Set(key, string(wd), timeOut...)
}
func (c *dbCache) SetNx(key, val string, timeOut ...int) error {
	rv := "-not-exitxs-key-"
	if c.Get(key, rv) == rv {
		return fmt.Errorf("%s 存在", key)
	}
	return c.Set(key, val, timeOut...)
}
func (c *dbCache) SetObjNx(key string, obj any, timeOut ...int) error {
	wd, e := json.Marshal(obj)
	if e != nil {
		return e
	}
	return c.SetNx(key, string(wd), timeOut...)
}
func (c *dbCache) Incr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.db.Exec(fmt.Sprintf("UPDATE %s SET val=val+%d WHERE key=?", c.tb, sv), key)
	if e != nil {
		c.Set(key, fmt.Sprintf("%d", sv))
		return sv
	} else {
		ll, e := r.RowsAffected()
		if e != nil {
			c.Set(key, fmt.Sprintf("%d", sv))
			return sv
		}
		if ll > 0 {
			r := c.Get(key, "0")
			ret, _ := strconv.Atoi(r)
			return ret
		} else {
			c.Set(key, fmt.Sprintf("%d", sv))
			return sv
		}
	}

}
func (c *dbCache) Decr(key string, step ...int) int {
	sv := 1
	if len(step) > 0 {
		sv = step[0]
	}
	r, e := c.db.Exec(fmt.Sprintf("UPDATE %s SET val=val-%d WHERE key=?", c.tb, sv), key)
	if e != nil {
		c.Set(key, fmt.Sprintf("%d", -sv))
		return sv
	} else {
		ll, e := r.RowsAffected()
		if e != nil {
			c.Set(key, fmt.Sprintf("%d", -sv))
			return sv
		}
		if ll > 0 {
			r := c.Get(key, "0")
			ret, _ := strconv.Atoi(r)
			return ret
		} else {
			c.Set(key, fmt.Sprintf("%d", -sv))
			return sv
		}
	}
}
func (c *dbCache) Del(key ...string) error {
	_, e := c.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE key IN (?)", c.tb), strings.Join(key, "','"))
	return e
}
func (c *dbCache) Flush() error {
	_, e := c.db.Exec(fmt.Sprintf("DELETE FROM %s", c.tb))
	return e
}
func (c *dbCache) Lock(key string, callfunc func(), timeOut ...int) error {
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
