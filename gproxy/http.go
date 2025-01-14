package gproxy

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

// http代理功能
func HttpProxy(con net.Conn) {
	buf := make([]byte, 4096)
	_, e := con.Read(buf)
	if e != nil {
		con.Close()
		return
	}

	var method, host, address string
	fmt.Sscanf(string(buf[:bytes.IndexByte(buf, '\n')]), "%s%s", &method, &host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err)
		return
	}

	if hostPortURL.Opaque == "443" { //https访问
		address = hostPortURL.Scheme + ":443"
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
		} else {
			address = hostPortURL.Host
		}
	}
	//获得了请求的host和port，就开始拨号吧
	client, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		fmt.Printf("链接请求")
		con.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	}
	NetCopy(con, client)
}
