package gapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ghf-go/glib/gutils"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

// 数据库Sesssion
func dbSession(conf gutils.ConfUrl) Handle {
	return func(c *Content) {

	}
}

// JWT session
func jwtSession(conf gutils.ConfUrl) Handle {
	jwtSignKey := []byte(conf.Get("key", "1234567890123456"))
	timeout := conf.GetInt("exipre", 1800)
	name := conf.Get("name", "Token")
	return func(c *Content) {
		token := sessionGetName(c, name)
		data := &jwt.RegisteredClaims{}
		if token != "" {
			t, e := jwt.ParseWithClaims(token, data, func(t *jwt.Token) (interface{}, error) {
				return jwtSignKey, nil
			})
			if e == nil && t.Valid { // && data.ExpiresAt != nil && data.ExpiresAt.Sub(time.Now()) > 0 {
				c.SetUserID(data.ID)
			}

		}
		c.Next()
		if c.IsLogin() {
			data.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(int64(timeout))))
			t := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
			if sendToken, e := t.SignedString(jwtSignKey); e == nil {
				sessionSetName(c, name, sendToken, timeout)
			}

		}
	}
}

// cache Session
func cacheSession(conf gutils.ConfUrl) Handle {
	ca := GetCache(conf.Host())
	timeout := conf.GetInt("exipre", 1800)
	name := conf.Get("name", "Token")
	return func(c *Content) {
		token := sessionGetName(c, name)
		if token != "" {
			c.SetUserID(ca.Get("uid", "0"))
		}
		c.Next()
		if c.IsLogin() {
			ca.Set(token, fmt.Sprintf("%d", c.GetUserID()), timeout)
			sessionSetName(c, name, token, timeout)
		}
	}
}

func sessionGetName(c *Content, name string) string {
	token := c.GetRequest().Header.Get(name)
	if token == "" {
		c, e := c.GetRequest().Cookie(name)
		if e == nil {
			token = c.Value
		}
	}
	if token == "" {
		token = c.GetRequest().URL.Query().Get(name)
	}
	if token == "" {
		token = uuid.NewV4().String()
	}
	return token
}
func sessionSetName(c *Content, name, val string, timeout int) {
	if c.IsLogin() {
		c.GetResponseWriter().Header().Add(name, val)
		http.SetCookie(c.GetResponseWriter(), &http.Cookie{
			Name:   name,
			Value:  val,
			MaxAge: timeout,
			Path:   "/",
		})
	}
}
