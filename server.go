package M

import (
	"fmt"
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
		fmt.Println("fail  to listen")
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
	fmt.Println(h.HandshakeOfWS(c))
	u := NewUnit(c)
	r_test.AddUnit(u)
	for {
		a, _ := u.Read()
		if a == nil {
			continue
		}
		r_test.BroadcastByUnit(u, a)
	}
}
