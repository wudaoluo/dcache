package service

import (
	"io"
	"net"
	"sync"

	"github.com/ivpusic/grpool"
	"github.com/wudaoluo/dcache/cache"
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/dcache/socket"
	"github.com/wudaoluo/golog"
)

//TODO 根据情况调整
const MAX_RECEIVE_DATA = 1024

type mux struct {
	datas        chan *internal.Data
	processDatas chan *internal.Data
	//maxReceiveData int32
	// notify a read event
	chReadEvent chan struct{}

	conn     socket.Socker
	dataLock sync.Mutex

	dataNotify chan struct{}
	connDead   chan struct{}
	worker     int
}

func NewMux(conn net.Conn) *mux {
	return &mux{
		datas:        make(chan *internal.Data, MAX_RECEIVE_DATA),
		processDatas: make(chan *internal.Data, MAX_RECEIVE_DATA),
		//maxReceiveData:MAX_RECEIVE_DATA,
		chReadEvent: make(chan struct{}),
		conn:        socket.NewTcpConn(conn),

		dataNotify: make(chan struct{}),
		connDead:   make(chan struct{}),
	}

}

func (m *mux) Operate() {
	defer close(m.processDatas)
	pool := grpool.NewPool(4, 96)
	defer pool.Release()
	for d := range m.datas {
		data := d
		pool.JobQueue <- func() {
			switch data.Op {
			case internal.OP_REQ_GET:
				value, ok := cache.Get(data.Key)
				if ok {
					data.Op = internal.OP_RESP_200
					data.Value = value
				} else {
					data.Op = internal.OP_RESP_404
					data.Value = []byte("key not found")
				}
			case internal.OP_REQ_PUT:
				cache.Set(data.Key, data.Value)

				data.Op = internal.OP_RESP_200
				data.Value = nil
			case internal.OP_REQ_DEL:
				cache.Del(data.Key)

				data.Op = internal.OP_RESP_200
				data.Value = nil
			default:
				data.Op = internal.OP_RESP_402
				data.Value = []byte("Reserved, used in the future")

			}

			m.processDatas <- data
		}

	}
}

func (m *mux) Read() {
	var err error

	defer func() {
		close(m.datas)
		m.conn.Close()
	}()
	for {
		var data = &internal.Data{}
		err = m.conn.ReadMsg(data)
		if err == io.EOF {
			golog.Info("mux.Read", "clientIP", m.conn.RemoteIP(), "err", "io.EOF")
			return
		}

		if err != nil {
			golog.Error("mux.Read", "clientIP", m.conn.RemoteIP(), "err", err)
			return
		}

		m.datas <- data

	}
}

func (m *mux) Write() {
	var err error
	for data := range m.processDatas {
		_, err = m.conn.WriteMsg(data)
		if err != nil {
			golog.Error("mux.Write", "err", err)
			m.conn.Close()
			return
		}
	}

}
