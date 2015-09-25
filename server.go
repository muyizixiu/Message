package M

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
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
	runtime.GOMAXPROCS(2)
	listener, err = net.Listen(proto, addr)
	http.HandleFunc("/", chatRoom)
	go http.ListenAndServe(":80", nil)
	if checkErr(err) {
		os.Exit(1)
	}
}
func chatRoom(w http.ResponseWriter, r *http.Request) {
	f, _ := os.Open("/root/Message/src/M/client.html")
	defer f.Close()
	io.Copy(w, f)
}
func Serve() {
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
	WriteHead(c, "HTTP/1.1 101 Switching Protocols\r\n")
	str, _ := GetHeadOfWS(c)
	WriteHead(c, str)
	EndHead(c)
	fmt.Println(c)
	u := NewUnit(c)
	r_test.AddUnit(u)
	for {
		a, err := u.Read()
		if err != nil {
			return
		}
		b, err := decode(a)
		if checkErr(err) {
			continue
		}
		ip := " <div class='ip'>" + c.RemoteAddr().String()[8:] + "</div>"
		b = append(b, []byte(ip)...)
		d, _ := encode(b)
		if a == nil {
			continue
		}
		r_test.BroadcastByUnit(u, d)
		fmt.Println(r_test.CountNumber())
	}
}
