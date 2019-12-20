package service

import (
	"github.com/wudaoluo/dcache/cache"
	"github.com/wudaoluo/dcache/plugin/session"
	"io"
	"runtime"
	"sync"

	"github.com/ivpusic/grpool"
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/golog"
)

//TODO 根据情况调整
const MAX_RECEIVE_DATA = 1024

//numWorkers和jobQueueLen 根据情况调整
var pool = grpool.NewPool(runtime.NumCPU()*2-1, 96)

type mux struct {
	datas        chan *internal.Data //读取client队列
	processDatas chan *internal.Data //发送client队列
	//maxReceiveData int32
	// notify a read event
	chReadEvent chan struct{}

	sess     *session.Session       //tcp 连接session
	dataLock sync.Mutex

	dataNotify chan struct{}
	connDead   chan struct{}
	worker     int
}

func NewMux(sess *session.Session) *mux {
	return &mux{
		datas:        make(chan *internal.Data, MAX_RECEIVE_DATA),
		processDatas: make(chan *internal.Data, MAX_RECEIVE_DATA),
		chReadEvent: make(chan struct{}),
		sess:sess,
		dataNotify: make(chan struct{}),
		connDead:   make(chan struct{}),
	}

}

func (m *mux) Operate() {
	defer func() {
		golog.Info("Operate close(m.processDatas)")
		close(m.processDatas)
	}()

	defer pool.Release()
	for d := range m.datas {
		data := d
		pool.JobQueue <- func() {
			value, ok := cache.Get(data.Key)
			if ok {
				data.Op = internal.OP_RESP_200
				data.Value = value
			} else {
				data.Op = internal.OP_RESP_404
				data.Value = []byte("key not found")
			}
			//对数据处理后写入发送client队列
			m.processDatas <- data
		}

	}
}

//读取数据写入 client读取队列
func (m *mux) Read() {
	var err error
	defer func() {
		golog.Info("(m *mux) Read() close(m.datas)")
		close(m.datas)
		session.HubClose(m.sess)
	}()
	for {
		var data = &internal.Data{}
		err = m.sess.ReadMsg(data)
		if err == io.EOF {
			golog.Info("mux.Read", "clientIP", m.sess.RemoteIP(), "err", "io.EOF")
			return
		}

		if err != nil {
			golog.Error("mux.Read", "clientIP", m.sess.RemoteIP(), "err", err)
			return
		}

		m.datas <- data
	}
}

//接收数据写入 client发送队列
func (m *mux) Write() {
	var err error
	for data := range m.processDatas {
		_, err = m.sess.WriteMsg(data)
		if err != nil {
			golog.Error("mux.Write", "err", err)
			session.HubClose(m.sess)
			return
		}
	}

}
