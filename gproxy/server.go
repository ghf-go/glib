package gproxy

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{}

// 新建http原生代理
func NewServerHttp(port int) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e == nil {
			go HttpProxy(c)
		}
	}
}

// 新建 sock原生代理
func NewServerSock5(port int) {
	l, e := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e == nil {
			go Sock5Proxy(c)
		}
	}
}

// 新建ws http代理
func NewServerHttpWs(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			return
		}
		ws, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go HttpWsProxy(ws)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// 新建 ws Sock代理
func NewServerSock5Ws(port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			return
		}
		ws, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		go Sock5WsProxy(ws)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
