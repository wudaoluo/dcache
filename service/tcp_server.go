package service

import (
	"net"
	"time"

	"github.com/wudaoluo/golog"

	"github.com/wudaoluo/dcache/socket"

	"io"

	"github.com/wudaoluo/dcache/internal"
)

type tcpServer struct {
	listen string
}

func NewTcpServer(listen string) Service {
	return &tcpServer{listen: listen}
}

func (t *tcpServer) Run() {
	ln, err := net.Listen("tcp", t.listen)

	if err != nil {
		golog.Error("tcpServer.Run", "err", err)
		return
	}

	defer ln.Close()
	golog.Info("tcpServer.Run", "tcp server listening addr:", t.listen)

	var tempDelay time.Duration
	for {
		conn, err := ln.Accept()
		//if err != nil {
		//	select {
		//	case <-ctx.Done():
		//		golog.Warn("tcpServer.run", "msg", "ctx Close tcp exit")
		//		return
		//	default:
		//
		//	}
		//}
		if err != nil {

			/*如果错误是暂时的,那么sleep一定时间在提供服务,否则就直接 return退出程序*/
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				golog.Warn("tcpServer.Run", "sleep time", tempDelay, "err", err)

				time.Sleep(tempDelay)
				continue
			} else {
				golog.Error("tcpServer.Run", "err", err)
				return
			}
		}

		go t.newHandler(conn)

	}
}

func (t *tcpServer) newHandler(conn net.Conn) {
	m := NewMux(conn)
	go m.Read()
	go m.Write()
	go m.Operate()
}

func (t *tcpServer) handler(conn net.Conn) {
	c := socket.NewTcpConn(conn)
	defer c.Close()
	var err error
	for {
		var data = &internal.Data{}
		err = c.ReadMsg(data)
		if err == io.EOF {
			golog.Info("tcpServer.handler", "clientIP", c.RemoteIP(), "err", "io.EOF")
			return
		}
		if err != nil {
			golog.Error("tcpServer.handler", "clientIP", c.RemoteIP(), "err", err)
			return
		}

		operate(data)

		_, err = c.WriteMsg(data)
		if err != nil {
			golog.Error("tcpServer.handler", "err", err)
			return
		}
	}

}
