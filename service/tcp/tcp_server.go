package tcp

import (
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/dcache/plugin/session"
	"github.com/wudaoluo/dcache/service"
	"github.com/wudaoluo/dcache/socket"
	"io"

	"net"
	"time"

	"github.com/wudaoluo/golog"
	"golang.org/x/net/netutil"
)

type tcpServer struct {
	listen  string
	maxConn int
	mux bool
}

func NewTcpServer(listen string, maxConn int,mux bool) service.Service {
	return &tcpServer{
		listen: listen,
		maxConn: maxConn,
		mux:mux,
	}
}

func (t *tcpServer) Run() {
	ln, err := net.Listen("tcp", t.listen)

	if err != nil {
		golog.Error("tcpServer.Run", "err", err)
		return
	}

	limitListener := netutil.LimitListener(ln, t.maxConn)

	defer ln.Close()
	golog.Info("tcpServer.Run", "tcp server listening addr:", t.listen, "maxconn", t.maxConn)

	var tempDelay time.Duration
	for {
		conn, err := limitListener.Accept()
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

		if t.mux {
			sess := session.NewSession(socket.NewTcpConn(conn))
			session.HubAdd(sess)
			go t.newHandler(sess)
		}else {
			go t.handler(conn)
		}

	}
}

func (t *tcpServer) newHandler(sess *session.Session) {
	m := service.NewMux(sess)
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

		service.Operate(data)

		_, err = c.WriteMsg(data)
		if err != nil {
			golog.Error("tcpServer.handler", "err", err)
			return
		}
	}

}
