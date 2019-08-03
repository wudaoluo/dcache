package internal

import (
	"github.com/spf13/cast"
)

const (
	PROJECT_NAME = "dcache"
	VERSION = "0.1"

)

type Port int32

const (
	TCP_PORT  Port = 7777
	GRPC_PORT Port = 7778

)

func (p Port) string() string {
	return cast.ToString(int32(p))
}

func (p Port) GetAddr(listen string) string {
	return listen + ":" + p.string()
}