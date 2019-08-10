package socket

import (
	"net"
	"io"
	"encoding/binary"
	"fmt"
	"github.com/wudaoluo/dcache/internal"
)


type tcpConn struct {
	net.Conn
	closed bool
}

func NewTcpConn(conn net.Conn) socker {
	return &tcpConn{
		conn,
		false,
	}
}


func (tc *tcpConn) Close() {
	if !tc.closed {
		tc.Close()
	}
}


func (tc *tcpConn) ReadMsg(req *internal.Req) error {
	msgHeadBuf := make([]byte,MSG_HEAD_LEN,MSG_HEAD_LEN)
	msgLen,err := tc.ReadLen(msgHeadBuf)
	if err != nil {
		return err
	}

	fmt.Println("msgLen",msgLen)
	//TODO ERR 统一处理
	if  msgLen > MAX_MSG_LEN {
		return fmt.Errorf("message too long")
	}

	if msgLen < MIN_MSG_LEN {
		return  fmt.Errorf("message too short")
	}

	msgBuf := make([]byte,msgLen)
	if _, err := io.ReadFull(tc, msgBuf); err != nil {
		return err
	}

	KeyLen := int(binary.BigEndian.Uint16(msgBuf[retainFront:keyHeadFront]))

	req.Op = msgBuf[:OP_LEN][0]
	fmt.Println(msgBuf[:OP_LEN][0])
	req.Retain = msgBuf[OP_LEN:retainFront][0]

	req.Key = msgBuf[keyHeadFront:keyHeadFront+KeyLen]
	if req.IsPut() {
		fmt.Println(valueHeadFront+KeyLen)
		//ValueLen := int(binary.BigEndian.Uint16(msgBuf[keyHeadFront+KeyLen:KeyLen+valueHeadFront]))
		req.Value = msgBuf[valueHeadFront+KeyLen:]
	}

	fmt.Println("op",int(req.Op),",key:",string(req.Key),",value:",string(req.Value),",re:",int(req.Retain))
	return nil
}

func (tc *tcpConn) ReadLen(b []byte) (n int,err error) {
	// read len
	if _, err := io.ReadFull(tc, b); err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint16(b)),nil
}


//TODO ERR 统一处理
func (tc *tcpConn) WriteMsg(op byte,key []byte,b []byte) (n int,err error) {
	msgLen := len(b)  + len(key)
	//
	if msgLen > MAX_MSG_LEN {
		return 0,fmt.Errorf("message too long")
	}

	if msgLen < MIN_MSG_LEN {
		return 0,fmt.Errorf("message too short")
	}
	//
	msg := make([]byte,MSG_HEAD_LEN,MSG_HEAD_LEN)
	binary.BigEndian.PutUint16(msg, uint16(msgLen+valueHeadFront))
	//
	keyLenBuf := make([]byte,KEY_HEAD_LEN)
	binary.BigEndian.PutUint16(keyLenBuf, uint16(len(key)))
	//
	valueLenBuf := make([]byte,VALUE_HEAD_LEN)
	binary.BigEndian.PutUint16(valueLenBuf, uint16(len(b)))
	//
	msg = append(msg,op,byte(0))
	msg = append(msg,keyLenBuf...)
	msg = append(msg,key...)
	msg = append(msg,valueLenBuf...)
	msg =  append(msg,b...)
	return tc.Write(msg)
}