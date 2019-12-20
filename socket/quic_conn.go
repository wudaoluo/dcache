package socket

import (
	"encoding/binary"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/golog"
	"io"
	"time"
)

type quicStream struct {
	readAge time.Duration
	stream quic.Stream
	closed bool
}

func NewQuicStream(stream quic.Stream) Socker {
	return &quicStream{
		stream:stream,
		closed:false,
		readAge:10*time.Second,  //连接默认存活时间
	}
}

func (qs *quicStream) RemoteIP() string {
	return ""
}

func (qs *quicStream) StreamId() quic.StreamID {  //int64
	return qs.stream.StreamID()
}

func (qs *quicStream) Close() {
	if !qs.closed {
		golog.Info("quicConn.Close", "streamid", qs.StreamId())
		_ = qs.stream.Close()
	}
}


func (qs *quicStream) ReadMsg(data *internal.Data) error {
	qs.stream.SetReadDeadline(time.Now().Add(qs.readAge))
	msgHeadBuf := make([]byte, MSG_HEAD_LEN, MSG_HEAD_LEN)
	msgLen, err := qs.ReadLen(msgHeadBuf)
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
	if _, err := io.ReadFull(qs.stream, msgBuf); err != nil {
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

func (qs *quicStream) ReadLen(b []byte) (n int, err error) {
	// read len
	if _, err := io.ReadFull(qs.stream, b); err != nil {
		return 0, err
	}

	return int(binary.BigEndian.Uint16(b)), nil
}

//TODO ERR 统一处理
func (qs *quicStream) WriteMsg(data *internal.Data) (n int, err error) {
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
	return qs.stream.Write(msg)
}