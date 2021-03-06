package socket

import (
	"encoding/binary"
	"fmt"
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/golog"
	"io"
	"net"
	"time"
)

type tcpConn struct {
	readAge time.Duration
	net.Conn
	closed bool
}

func NewTcpConn(conn net.Conn) Socker {
	return &tcpConn{
		Conn:conn,
		closed:false,
		readAge:10*time.Second,  //连接默认存活时间
	}
}

func (tc *tcpConn) Close() {
	if !tc.closed {
		golog.Info("tcpConn.Close", "clientIP", tc.RemoteAddr())
		tc.Conn.Close()
	}
}

func (tc *tcpConn) RemoteIP() string {
	return tc.Conn.RemoteAddr().String()
}

func (tc *tcpConn) ReadMsg(data *internal.Data) error {
	tc.Conn.SetReadDeadline(time.Now().Add(tc.readAge))
	msgHeadBuf := make([]byte, MSG_HEAD_LEN, MSG_HEAD_LEN)
	msgLen, err := tc.ReadLen(msgHeadBuf)
	if err != nil {
		return err
	}

	if msgLen > MAX_MSG_LEN {
		return fmt.Errorf("message too long %d", msgLen)
	}

	if msgLen < MIN_MSG_LEN {
		return fmt.Errorf("message too short")
	}

	msgBuf := make([]byte, msgLen)
	if _, err := io.ReadFull(tc, msgBuf); err != nil {
		return err
	}

	KeyLen := int(binary.BigEndian.Uint16(msgBuf[retainFront:keyHeadFront]))

	data.Op = msgBuf[:OP_LEN][0]
	data.Retain = msgBuf[OP_LEN:retainFront][0]

	data.Key = msgBuf[keyHeadFront : keyHeadFront+KeyLen]
	if data.IsValue() {
		data.Value = msgBuf[valueHeadFront+KeyLen:]
	}
	return nil
}

func (tc *tcpConn) ReadLen(b []byte) (n int, err error) {
	// read len
	if _, err := io.ReadFull(tc, b); err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint16(b)), nil
}

//TODO ERR 统一处理
func (tc *tcpConn) WriteMsg(data *internal.Data) (n int, err error) {
	msgLen := len(data.Value) + len(data.Key) + valueHeadFront

	if msgLen > MAX_MSG_LEN {
		return 0, fmt.Errorf("message too long")
	}

	if msgLen < MIN_MSG_LEN {
		return 0, fmt.Errorf("message too short")
	}
	//
	msg := make([]byte, MSG_HEAD_LEN, MSG_HEAD_LEN)
	binary.BigEndian.PutUint16(msg, uint16(msgLen))
	//
	keyLenBuf := make([]byte, KEY_HEAD_LEN)
	binary.BigEndian.PutUint16(keyLenBuf, uint16(len(data.Key)))
	//
	valueLenBuf := make([]byte, VALUE_HEAD_LEN)
	binary.BigEndian.PutUint16(valueLenBuf, uint16(len(data.Value)))
	//
	msg = append(msg, data.Op, data.Retain)
	msg = append(msg, keyLenBuf...)
	msg = append(msg, data.Key...)
	msg = append(msg, valueLenBuf...)
	msg = append(msg, data.Value...)
	return tc.Write(msg)
}
