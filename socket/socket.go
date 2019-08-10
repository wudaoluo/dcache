package socket

import "github.com/wudaoluo/dcache/internal"

const (
	MIN_MSG_LEN    int = 4
	MAX_MSG_LEN    int = 4096
	MSG_HEAD_LEN   int = 2
	OP_LEN         int = 1 //操作位
	RETAIN_LEN         = 1 //保留长度
	KEY_HEAD_LEN   int = 2
	VALUE_HEAD_LEN int = 2
)

var (
	retainFront    = OP_LEN + RETAIN_LEN           //2
	keyHeadFront   = retainFront + KEY_HEAD_LEN    //4
	valueHeadFront = keyHeadFront + VALUE_HEAD_LEN // 6

)

type Socker interface {
	ReadMsg(d *internal.Data) error
	WriteMsg(d *internal.Data) (n int, err error)
	RemoteIP() string
	Close()
}
