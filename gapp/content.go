package gapp

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Content struct {
	ReqID         string //请求id
	r             *http.Request
	w             http.ResponseWriter
	handles       []Handle
	reqIP         string
	isAbort       bool
	currentNext   int
	responseBytes []byte
	userId        uint64
	ctx           context.Context
	clientlang    string
	tpl           *template.Template
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 获取客户端IP
func (c *Content) GetIP() string {
	if c.r == nil {
		return "127.0.0.1"
	}
	if c.reqIP == "" {
		ret := c.r.Header.Get("ipv4")
		if ret != "" {
			c.reqIP = ret
			return c.reqIP
		}
		ret = c.r.Header.Get("X-Forwarded-For")
		if ret != "" {
			rs := strings.Split(ret, ",")
			if rs[0] != "" {
				c.reqIP = rs[0]
				return c.reqIP
			}
		}
		ret = c.r.Header.Get("XForwardedFor")
		if ret != "" {
			rs := strings.Split(ret, ",")
			if rs[0] != "" {
				c.reqIP = rs[0]
				return c.reqIP
			}
		}
		ret = c.r.Header.Get("X-Real-Ip")
		if ret != "" {
			rs := strings.Split(ret, ",")
			if rs[0] != "" {
				c.reqIP = rs[0]
				return c.reqIP
			}
		}
		ret = c.r.Header.Get("X-Real-IP")
		if ret != "" {
			rs := strings.Split(ret, ",")
			if rs[0] != "" {
				c.reqIP = rs[0]
				return c.reqIP
			}
		}
		ret = c.r.RemoteAddr
		if ret != "" {
			ips := strings.Split(ret, ":")
			c.reqIP = ips[0]
			return c.reqIP
		}

		c.reqIP = "unknow"
	}
	return c.reqIP
}

func (c *Content) Abort() {
	c.isAbort = true
}
func (c *Content) Next() {
	if c.r == nil {
		return
	}
	if c.currentNext < len(c.handles) {
		ci := c.currentNext
		c.currentNext++
		if !c.isAbort {
			c.handles[ci](c)
		}
	}

}

// 是否开启DEBUG 日志输出模式
func (c *Content) IsDebug() bool {
	return config.App.Debug
}

// 国际化语言
func (c *Content) Lang(key string) string {
	if c.clientlang == "" {
		kks := strings.Split(c.r.Header.Get("accept-language"), ",")
		for _, item := range kks {
			iks := strings.Split(item, ";")
			for _, v := range iks {
				if strings.Index(v, "-") > 0 {
					c.clientlang = strings.ToLower(v)
					return config.Lang.Lang(key, c.clientlang)
				}
			}
		}
	}
	return config.Lang.Lang(key, c.clientlang)
}

// 绑定数据
func (c *Content) BindJson(obj any) error {
	if c.r == nil {
		return errors.New("计划任务不能解析参数")
	}
	body := c.r.Body
	defer body.Close()
	data, e := io.ReadAll(body)
	if e != nil {
		return e
	}
	return json.Unmarshal(data, obj)
}

// 绑定提交的xml数据
func (c *Content) BindXml(obj any) error {
	if c.r == nil {
		return errors.New("计划任务不能解析参数")
	}
	body := c.r.Body
	defer body.Close()
	data, e := io.ReadAll(body)
	if e != nil {
		return e
	}
	return xml.Unmarshal(data, obj)
}

// 接口正常返回
func (c *Content) SuccessJson(data any) {
	c.json(200, "", data)
}

// 接口保存信息
func (c *Content) FailJson(code int, msg string) {
	c.json(code, msg, nil)
}

// 返回 xml
func (c *Content) Xml(data string) {
	if c.r == nil {
		return
	}
	c.w.Header().Set("content-type", "text/xml;charset=utf8")
	c.responseBytes = []byte(data)

}

// 返回 html内容
func (c *Content) Out(data string) {
	if c.r == nil {
		return
	}
	c.responseBytes = []byte(data)
}

// 输出JSON信息
func (c *Content) json(code int, msg string, data any) {
	if c.r == nil {
		return
	}
	c.w.Header().Set("content-type", "application/json;charset=utf8")
	ret := map[string]any{
		"code": code,
		"msg":  msg,
		"data": data,
	}
	dd, e := json.Marshal(ret)
	if e != nil {
		panic(e.Error())
	}
	c.responseBytes = dd

}

// 刷新缓存
func (c *Content) flush() {
	if c.r == nil {
		return
	}
	c.w.(http.Flusher).Flush()
}

// 开启Event事件
func (c *Content) Sse(call func(s *Sse)) {
	if c.r == nil {
		return
	}
	c.w.Header().Set("Content-Type", "text/event-stream")
	// 这行代码设置 HTTP 响应的 Cache-Control 为 no-cache，告诉浏览器不要缓存此响应。
	c.w.Header().Set("Cache-Control", "no-cache")
	// 这行代码设置 HTTP 响应的 Connection 为 keep-alive，保持长连接，以便服务器可以持续发送事件到客户端。
	c.w.Header().Set("Connection", "keep-alive")
	// 这行代码设置 HTTP 响应的自定义头部 X-Accel-Buffering 为 no，用于禁用某些代理或 Web 服务器（如 Nginx）的缓冲。这有助于确保服务器发送事件在传输过程中不会受到缓冲影响
	c.w.Header().Set("X-Accel-Buffering", "no")
	c.w.Header().Set("Access-Control-Allow-Origin", "*")
	c.flush()
	call(&Sse{c: c, isClose: false, key: c.ReqID})
}

// 开始websocket
func (c *Content) WebSocket(call func(con *websocket.Conn)) {
	if c.r == nil {
		return
	}
	conn, err := upgrader.Upgrade(c.w, c.r, nil)
	if err != nil {
		fmt.Printf("链接失败 %s\n", err.Error())
		log.Println(err)
		return
	}
	call(conn)
	defer fmt.Printf("链接管理了\n")
	defer conn.Close()

	// if c.r.Header.Get("Upgrade") != "websocket" {
	// 	return errors.New("协议错误")
	// }
	// if c.r.Header.Get("Sec-WebSocket-Version") != "13" {
	// 	return errors.New("协议错误")
	// }
	// k := c.r.Header.Get("Sec-WebSocket-Key")
	// if k == "" {
	// 	return errors.New("协议错误")
	// }
	// c.w.WriteHeader(http.StatusSwitchingProtocols)
	// c.w.Header().Set("Upgrade", "websocket")
	// c.w.Header().Set("Sec-WebSocket-Version", "13")
	// c.w.Header().Set("Connection", "Upgrade")
	// dd := sha1.Sum([]byte(k + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	// c.w.Header().Set("Sec-WebSocket-Accept", base64.StdEncoding.EncodeToString(dd[:]))
	// c.flush()
	// return nil
}

// 显示模版
func (c *Content) Display(tpl string, data any) {
	c.responseBytes = []byte(c.Template(tpl, data))
}

// 获取模版信息
func (c *Content) Template(tpl string, data any) string {
	buf := &bytes.Buffer{}
	if e := c.tpl.ExecuteTemplate(buf, tpl, data); e != nil {
		return ""
	}
	return buf.String()
}

// 显示模版
func (c *Content) DisplayLayout(layout, tpl string, data any) {
	content := c.Template(tpl, data)
	out := c.Template(layout, data)
	out = strings.ReplaceAll(out, "__LAYOUT_CONTENT__", content)
	c.responseBytes = []byte(out)
}

// 获取配置信息
func (c *Content) GetConf() *conf {
	return config
}
func (c *Content) Flush() {
	if c.r == nil {
		debug(" %s", string(c.responseBytes))
		return
	}
	c.w.Write(c.responseBytes)
}
func (c *Content) GetRequest() *http.Request {
	if c.r == nil {
		return nil
	}
	return c.r
}
func (c *Content) GetResponseWriter() http.ResponseWriter {
	if c.r == nil {
		return nil
	}
	return c.w
}
func (c *Content) SetUserID(uid string) {
	if c.r == nil {
		return
	}
	if id, e := strconv.ParseUint(uid, 10, 64); e == nil {
		c.userId = id
	}
}

func (c *Content) GetUserID() uint64 {
	if c.r == nil {
		return 0
	}
	return c.userId
}
func (c *Content) IsLogin() bool {
	if c.r == nil {
		return false
	}
	return c.userId > 0
}
func (c *Content) GetContext() context.Context {
	return c.ctx
}
