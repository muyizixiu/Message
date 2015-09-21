package M

import (
	"fmt"
	"net"
	"sync"
)

func watch(v interface{}, line int) {
	fmt.Println(line, "  ", v)
}

const (
	MaxConn = 1000
)

var (
	unitPool struct {
		pool []*Unit
		mu   *sync.Mutex
	}
)

type Unit struct {
	conn *net.TCPConn
	id   int
}

func NewUnit(c *net.TCPConn) *unit {
	a := new(Unit{conn: c})
	unitPool.mu.Lock()
	unitPool = append(unitPool.pool, a)
	unitPool.mu.Unlock()
	a.id = len(unitPool) - 1
	return a
}
func (c *Unit) Close() {
	var head []*Unit
	var tail []*Unit
	unitPool.mu.Lock()
	head = unitPool.pool[:c.id]
	tail = unitPool.pool[c.id+1:]
	unitPool = append(head, tail...)
	unitPool.mu.Unlock()
	c.conn.Close()
}
func (c *Unit) Read() ([]byte, error) {
	b := make([]byte, 1024)
	var result []byte
	for {
		n, err := c.conn.Read(b)
		result = append(result, b[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			if checkErr(err) {
				return nil, err
			}

		}
	}
	return result, nil
}
func checkErr(err error) bool {
	if err == nil {
		return false
	}
	return true
}

type Room struct {
	players []*Unit
}
