package gapp

import "github.com/ghf-go/glib/gutils"

// 验证是否登录的中间件
func CheckoutLogin(c *Content) {
	if c.IsLogin() {
		c.Next()
	} else {
		c.FailJson(303, "账号没有登录")
	}
}

func SessionHandel(name ...string) Handle {
	kn := "default"
	if len(name) > 0 {
		kn = name[0]
	}
	if c, ok := config.App.Session[kn]; ok {
		cc := gutils.NewConfUrl(c)
		switch cc.Scheme() {
		case "cache":
			return cacheSession(cc)
		case "jwt":
			return jwtSession(cc)
		case "db":
			return dbSession(cc)
		}
	}
	return func(c *Content) {
		c.Next()
	}

}
