package gproxy

import (
	"fmt"
	"net"
)

// sock5代理功能
func Sock5Proxy(con net.Conn) {
	h1 := make([]byte, 3)
	_, e := con.Read(h1)
	if e != nil {
		// pserver.Debug("1 sock init err %s", e.Error())
		return
	}

	_, e = con.Write([]byte{0x05, 0x00})
	if e != nil {
		// pserver.Debug("2 sock init err %s", e.Error())
		return
	}

	h2 := make([]byte, 4)
	_, e = con.Read(h2)
	if e != nil {
		// pserver.Debug("3 sock init err %s", e.Error())
		return
	}

	addr := ""
	_, e = con.Write([]byte{0x05, 0x00, 0x00, h2[3]})
	if e != nil {
		// pserver.Debug("4 sock init err %s", e.Error())
		return
	}

	switch h2[3] {
	case 0x01:
		v := make([]byte, 4)
		_, e = con.Read(v)
		if e != nil {
			// pserver.Debug("5 sock init err %s", e.Error())
			return
		}

		addr = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
		fmt.Printf("链接地址 %s\n", addr)
		_, e = con.Write(v)
		if e != nil {
			// pserver.Debug("6 sock init err %s", e.Error())
			return
		}

	case 0x04:
		v := make([]byte, 16)
		_, e = con.Read(v)
		if e != nil {
			// pserver.Debug("7 sock init err %s", e.Error())
			return
		}

		_, e = con.Write(v)
		if e != nil {
			// pserver.Debug("8 sock init err %s", e.Error())
			return
		}
	case 0x03:
		l := make([]byte, 1)
		_, e = con.Read(l)
		if e != nil {
			// pserver.Debug("9 sock init err %s", e.Error())
			return
		}
		v := make([]byte, l[0])
		_, e = con.Read(v)
		if e != nil {
			// pserver.Debug("10 sock init err %s", e.Error())
			return
		}

		_, e = con.Write(l)
		if e != nil {
			// pserver.Debug("11 sock init err %s", e.Error())
			return
		}

		_, e = con.Write(v)
		if e != nil {
			// pserver.Debug("12 sock init err %s", e.Error())
			return
		}
		addr = string(v)
		fmt.Printf("链接地址 %s\n", addr)
	}
	if addr == "" {
		con.Close()
	}
	p := make([]byte, 2)
	_, e = con.Read(p)
	if e != nil {
		// pserver.Debug("13 sock init err %s", e.Error())
		return
	}

	port := (int(p[0]) << 8) + int(p[1])
	fmt.Printf("链接端口 %d - %d %d %s\n", port, int(p[0])<<8, p[1], string(p))
	desc, e := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if e != nil {
		// pserver.Debug("14 sock init err %s", e.Error())
		return
	}

	_, e = con.Write(p)
	if e != nil {
		// pserver.Debug("15 sock init err %s", e.Error())
		return
	}
	NetCopy(con, desc)
}
