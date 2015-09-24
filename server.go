package M

import (
	"handshake/h"
	"net"
	"os"
)

const (
	addr  = ":8080"
	proto = "tcp"
)

var (
	listener net.Listener
)

func init() {
	var err error
	listener, err = net.Listen(proto, addr)
	if checkErr(err) {
		os.Exit(1)
	}
}
func serve() {
	for {
		c, err := listener.Accept()
		if checkErr(err) {
			continue
		}
		go handle(c)
	}
}

var r_test = NewRoom()

func handle(c net.Conn) {
	u := NewUnit(c)
	r_test.AddUnit(u)
	for {
		a, _ := u.Read()
		b, _ := decode(a)
		d, _ := encode(b)
		if a == nil {
			continue
		}
		r_test.BroadcastByUnit(u, d)
	}
}
