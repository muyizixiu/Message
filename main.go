package M

import (
	"errors"
	"fmt"
	"net"
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

type Unit struct { //basic unit for each connection
	conn *net.TCPConn
	id   int
}

func NewUnit(c *net.TCPConn) *unit {
	a := new(Unit{conn: c})
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
func (c *Unit) Write(v []byte) {
	_, e := c.conn.Write(v)
	checkErr(e)
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
	tmp := new(Room{id: getRoomId()})
	room = append(room, tmp)
	return tmp
}
func (r Room) AddUnit(u *Unit) {
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
func (r Room) BroadcastByUnit(u *unit, v []byte) {
	for _, v := range r.players {
		if u == v {
			continue
		}
		v.Write([]byte)
	}
}
