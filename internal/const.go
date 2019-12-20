package internal

import (
	"github.com/spf13/cast"
)

const (
	PROJECT_NAME = "dcache"
	VERSION      = "0.1"
)

type Port int32

const (
	TCP_PORT  Port = 7777    //tcp
	QUIC_PORT  Port = 7777   //udp
	GRPC_PORT Port = 7778
)

func (p Port) string() string {
	return cast.ToString(int32(p))
}

func (p Port) GetAddr(listen string) string {
	return listen + ":" + p.string()
}

/*
1**	信息，服务器收到请求，需要请求者继续执行操作
2**	成功，操作被成功接收并处理
3**	重定向，需要进一步的操作以完成请求
4**	客户端错误，请求包含语法错误或无法完成请求
5**	服务器错误，服务器在处理请求的过程中发生了错误

https://www.runoob.com/http/http-status-codes.html
*/
const (
	//保留 = byte(0)
	OP_REQ_GET = byte(1)
	OP_REQ_PUT = byte(2)
	OP_REQ_DEL = byte(3)

	//http response code
	//保留 = byte(0)
	OP_RESP_200 = byte(1)
	OP_RESP_404 = byte(2)
	OP_RESP_302 = byte(3)
	OP_RESP_402 = byte(4) //保留，将来使用
)
