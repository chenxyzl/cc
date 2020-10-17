package session

import (
	"net"
	"time"
)


//IConn 通用连接
type IConn interface {
	Read(p []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	SetReadDeadline(t time.Time) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
}

type Message struct {
	Msg string `json:"msg"`
}

type OutMessage struct {
}

type Sess struct {
	conn IConn
}

func (s *Sess) ReadMessage() {

}
func (s *Sess) WriteMessage() {

}
