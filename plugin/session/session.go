package session

import (
	"github.com/wudaoluo/dcache/socket"
	"time"
)

type Session struct {
	socket.Socker
	Expiration int64
}

func NewSession(socket socket.Socker) *Session {
	return &Session{
		socket, time.Now().Add(30*time.Second).Unix(), //默认保留30秒 心跳时间 ( (30 -1)/3 )
}
}

func (s *Session) ID() string {
	return s.RemoteIP()
}
