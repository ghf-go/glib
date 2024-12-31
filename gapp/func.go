package gapp

import (
	"fmt"

	"github.com/ghf-go/glib/gcache"
	"github.com/ghf-go/glib/gutils"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 获取Cache
func GetCache(name ...string) gcache.Cache {
	cname := "default"
	if len(name) > 0 {
		cname = name[0]
	}
	if r, ok := mapCache[cname]; ok {
		return r
	}
	var cache gcache.Cache
	if c, ok := config.Caches[cname]; ok {
		cc := gutils.NewConfUrl(c)
		switch cc.Scheme() {
		case "redis":
			cache = gcache.NewRedisCache(cc)
			break
		case "memcache":
			cache = gcache.NewMemcache(cc)
			break
		case "db":
			cache = gcache.NewDbCache(GetDB(cc.Host()))
			break
		default:
			cache = gcache.NewFileCache(cc)
			break
		}
	}
	if cache != nil {
		mapCache[cname] = cache
	}
	return cache
}

// 获取Redis
func GetCacheRedis(name ...string) *redis.Client {
	cname := "default"
	if len(name) > 0 {
		cname = name[0]
	}
	if r, ok := mapRedis[cname]; ok {
		return r
	}
	if c, ok := config.Caches[cname]; ok {
		cc := gutils.NewConfUrl(c)
		ret := gcache.NewRedisClient(cc)
		mapRedis[cname] = ret
		return ret
	}
	return nil
}

// 获取数据库链接
func GetDB(name ...string) *gorm.DB {
	cname := "default"
	if len(name) > 0 {
		cname = name[0]
	}
	if r, ok := mapDb[cname]; ok {
		return r
	}
	if c, ok := config.Dbs[cname]; ok {
		cc := gutils.NewConfUrl(c)
		db, e := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cc.User(), cc.Pass(), cc.Host(), cc.Path())), &gorm.Config{})
		if e != nil {
			panic(e.Error())
		}
		db.Use(dbresolver.Register(dbresolver.Config{}).SetMaxIdleConns(cc.GetInt("MaxIdleConns", 3)).SetMaxOpenConns(cc.GetInt("MaxOpenConns", 5)).SetConnMaxIdleTime(cc.GetDuration("ConnMaxIdleTime", 1800)).SetConnMaxLifetime(cc.GetDuration("ConnMaxLifetime", 1800)))
		mapDb[cname] = db
		return db
	}

	return nil
}

// 处理事务
func Tx(db *gorm.DB, call func(tx *gorm.DB) (error, any)) (error, any) {
	tx := db.Begin()
	e, ret := call(tx)
	if e != nil {
		tx.Rollback()
		return e, ret
	}
	tx.Commit()
	return e, ret
}

// // 通过本机发送邮件
// func (c *Content) SendLocalMail(conname, to, subject string, isHtml bool, msg []byte) error {
// 	i := strings.Index(to, "@")
// 	host := to[i+1:]
// 	if sc, ok := c.confData.Stmp[conname]; ok {
// 		if dd, e := net.LookupMX(host); e == nil {
// 			content_type := ""
// 			if isHtml {
// 				content_type = "Content-Type: text/html; charset=UTF-8"
// 			} else {
// 				content_type = "Content-Type: text/plain" + "; charset=UTF-8"
// 			}
// 			msg = []byte("To: " + to + "\r\nFrom: " + sc.UserName + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + string(msg))

// 			return smtp.SendMail(dd[0].Host+":25", nil, sc.UserName, []string{to}, msg)
// 		}
// 		return errors.New("获取信息失败")
// 	}
// 	return errors.New("配置不存在" + conname)

// }

// // 发送邮件
// func (c *Content) SendMail(conname, to, subject string, isHtml bool, msg []byte) error {
// 	if sc, ok := c.confData.Stmp[conname]; ok {
// 		var auth smtp.Auth
// 		switch strings.ToUpper(sc.AuthType) {
// 		case "CRAMMD5":
// 			auth = smtp.CRAMMD5Auth(sc.UserName, sc.Passwd)
// 		case "HOTMAIL":
// 			auth = conf.NewHotmailStmpAuth(sc.UserName, sc.Passwd)
// 		default:
// 			auth = smtp.PlainAuth("", sc.UserName, sc.Passwd, sc.Host)
// 		}
// 		content_type := ""
// 		if isHtml {
// 			content_type = "Content-Type: text/html; charset=UTF-8"
// 		} else {
// 			content_type = "Content-Type: text/plain" + "; charset=UTF-8"
// 		}
// 		msg = []byte("To: " + to + "\r\nFrom: " + sc.UserName + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + string(msg))

// 		return smtp.SendMail(sc.Host, auth, sc.UserName, []string{to}, msg)
// 	}
// 	return errors.New("配置不存在")
// }
