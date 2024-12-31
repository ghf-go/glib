package gapp

import (
	"github.com/ghf-go/glib/gcache"
	"github.com/ghf-go/glib/gutils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 获取Cache
func GetCache(name ...string) gcache.Cache {
	cname := "default"
	if len(name) > 0 {
		cname = name[0]
	}
	if c, ok := config.Caches[cname]; ok {
		cc := gutils.NewConfUrl(c)
	}
	sysError("配置错误")
	return nil
}

// 获取Redis
func GetCacheRedis(conname ...string) *redis.Client {
	return nil
}

// 获取数据库链接
func GetDB(conname ...string) *gorm.DB {
	return nil
}

// 获取数据库
func (c *Content) GetDB(dbname ...string) *gorm.DB {
	conName := "default"
	if len(dbname) > 0 {
		conName = dbname[0]
	}
	if r, ok := dbCon[conName]; ok {
		return r
	}
	if dbconf, ok := c.confData.Dbs[conName]; ok {
		db, e := gorm.Open(mysql.Open(dbconf.Write), &gorm.Config{})
		if e != nil {
			panic(e.Error())
		}
		rs := []gorm.Dialector{}
		for _, rc := range dbconf.Reads {
			rs = append(rs, mysql.Open(rc))
		}
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(dbconf.Write)},
			Replicas: rs,
		}).SetMaxIdleConns(dbconf.ConMaxIdleTime).SetMaxOpenConns(dbconf.MaxOpenCons).SetConnMaxIdleTime(time.Minute * time.Duration(dbconf.ConMaxIdleTime)).SetConnMaxLifetime(time.Minute * time.Duration(dbconf.ConMaxLifeTime)))
		if c.confData.App.Debug {
			db = db.Debug()
		}
		dbCon[conName] = db
		return db
	} else {
		panic("数据配置错误 " + conName)
	}

}

// 获取缓存配置
func (c *Content) GetCache(conname ...string) *redis.Client {
	conName := "default"
	if len(conname) > 0 {
		conName = conname[0]
	}
	if r, ok := cacheCon[conName]; ok {
		return r
	}
	if rconf, ok := c.confData.Cache[conName]; ok {
		r := redis.NewClient(&redis.Options{
			Addr:            rconf.Host,
			Username:        rconf.UserName,
			Password:        rconf.Passwd,
			MinIdleConns:    rconf.MinIdleConns,
			MaxIdleConns:    rconf.MaxIdleConns,
			MaxActiveConns:  rconf.MaxActiveConns,
			ConnMaxIdleTime: time.Minute * time.Duration(rconf.ConnMaxIdleTime),
			ConnMaxLifetime: time.Minute * time.Duration(rconf.ConnMaxLifetime),
		})
		cacheCon[conName] = r
		return r
	}
	panic("缓存配置不存在" + conName)

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

// 通过本机发送邮件
func (c *Content) SendLocalMail(conname, to, subject string, isHtml bool, msg []byte) error {
	i := strings.Index(to, "@")
	host := to[i+1:]
	if sc, ok := c.confData.Stmp[conname]; ok {
		if dd, e := net.LookupMX(host); e == nil {
			content_type := ""
			if isHtml {
				content_type = "Content-Type: text/html; charset=UTF-8"
			} else {
				content_type = "Content-Type: text/plain" + "; charset=UTF-8"
			}
			msg = []byte("To: " + to + "\r\nFrom: " + sc.UserName + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + string(msg))

			return smtp.SendMail(dd[0].Host+":25", nil, sc.UserName, []string{to}, msg)
		}
		return errors.New("获取信息失败")
	}
	fmt.Println(c.confData.Stmp)
	return errors.New("配置不存在" + conname)

}

// 发送邮件
func (c *Content) SendMail(conname, to, subject string, isHtml bool, msg []byte) error {
	if sc, ok := c.confData.Stmp[conname]; ok {
		var auth smtp.Auth
		switch strings.ToUpper(sc.AuthType) {
		case "CRAMMD5":
			auth = smtp.CRAMMD5Auth(sc.UserName, sc.Passwd)
		case "HOTMAIL":
			auth = conf.NewHotmailStmpAuth(sc.UserName, sc.Passwd)
		default:
			auth = smtp.PlainAuth("", sc.UserName, sc.Passwd, sc.Host)
		}
		content_type := ""
		if isHtml {
			content_type = "Content-Type: text/html; charset=UTF-8"
		} else {
			content_type = "Content-Type: text/plain" + "; charset=UTF-8"
		}
		msg = []byte("To: " + to + "\r\nFrom: " + sc.UserName + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + string(msg))

		return smtp.SendMail(sc.Host, auth, sc.UserName, []string{to}, msg)
	}
	return errors.New("配置不存在")
}
