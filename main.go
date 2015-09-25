package M

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"
)

func watch(v interface{}, line int) {
	fmt.Println(line, "  ", v)
}

const (
	MaxConn = 1000
)

var (
	unitPool struct { //manage long connections
		pool    []*Unit //store connection
		idlePos []int   //record position of closed connection
		mu      *sync.Mutex
	}
)

func CountConnNumber() int {
	var i int
	fmt.Println(unitPool.pool)
	for _, v := range unitPool.pool {
		if v != nil {
			i++
		}
	}
	return i
}

func init() {
	runtime.GOMAXPROCS(2)
	unitPool.mu = new(sync.Mutex)
}

type Unit struct { //basic unit for each connection
	conn net.Conn
	id   int
}

func NewUnit(c net.Conn) *Unit {
	a := &Unit{conn: c}
	unitPool.mu.Lock()
	if len(unitPool.idlePos) < 10 {
		a.id = len(unitPool.pool)
		unitPool.pool = append(unitPool.pool, a)
	} else {
		a.id = unitPool.idlePos[0]
		unitPool.pool[a.id] = a
		unitPool.idlePos = unitPool.idlePos[1:]
	}
	unitPool.mu.Unlock()
	return a
}
func (c *Unit) Close() {
	unitPool.mu.Lock()
	unitPool.idlePos = append(unitPool.idlePos, c.id)
	unitPool.pool[c.id] = nil
	unitPool.mu.Unlock()
	c.conn.Close()
}
func (c *Unit) Read() ([]byte, error) {
	b := make([]byte, 1024)
	var result []byte
	for {
		n, err := c.conn.Read(b)
		fmt.Println(err)
		result = append(result, b[:n]...)
		if checkErr(err) {
			c.Close()
			return nil, err
		}

		if n < 1024 {
			break
		}
	}
	return result, nil
}
func (c *Unit) Write(v []byte) error {
	_, e := c.conn.Write(v)
	if checkErr(e) {
		c.Close()
		return e
	}
	return nil
}
func checkErr(err error) bool {
	if err == nil {
		return false
	}
	return true
}

var (
	roomId   int     //for a room to get a id number
	room     []*Room //store all Room
	idleRoom []int   //record closed Room
)

type Room struct {
	players []*Unit
	id      int
}

func getRoomId() int {
	if len(idleRoom) > 2 {
		a := idleRoom[0]
		idleRoom = idleRoom[1:]
		return a
	}
	roomId++
	return roomId
}
func NewRoom() *Room {
	tmp := &Room{id: getRoomId()}
	room = append(room, tmp)
	return tmp
}
func (r *Room) AddUnit(u *Unit) {
	r.players = append(r.players, u)
}
func (r Room) deleteUnit(u *Unit) {
	var tmp []*Unit
	for _, v := range r.players {
		if u == v {
			continue
		}
		tmp = append(tmp, v)
	}
	r.players = tmp
}
func GetRoom(id int) *Room {
	for _, v := range room {
		if id == v.id {
			return v
		}
	}
	checkErr(errors.New("there no room" + strconv.Itoa(id)))
	return nil
}
func (r Room) BroadcastByUnit(u *Unit, val []byte) {
	for i, v := range r.players {
		if v == nil {
			continue
		}
		if u == v {
			continue
		}
		err := v.Write(val)
		if err != nil {
			r.players[i] = nil
		}
	}
}
func (r Room) CountNumber() int {
	var i int
	for _, v := range r.players {
		if v != nil {
			i++
		}
	}
	return i
}
