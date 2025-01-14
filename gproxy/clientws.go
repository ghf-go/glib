package gproxy

import (
	"fmt"
	"net"

	"github.com/gorilla/websocket"
)

// 创建ws代理客户端
func NewClientProxy(port int, wsurl string) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	// pserver.Debug("开始监听服务")
	if e != nil {
		// pserver.Debug("服务启动异常 %s", e.Error())
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			// pserver.Debug("建立链接异常 %s", e.Error())
			continue
		}
		// pserver.Debug("clientIP %s", c.RemoteAddr().String())
		// go runClient(c)
		go func(con net.Conn, host string) {
			c, _, err := websocket.DefaultDialer.Dial(host, nil)
			if err != nil {
				// pserver.Debug("建立ws失败 %s", err.Error())
				return
			}
			netCopyWsCon(con, c)
		}(c, wsurl)
	}
}
