package service

import (
	"github.com/wudaoluo/dcache/cache"
	"github.com/wudaoluo/dcache/internal"
)

func Operate(data *internal.Data) {
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

}
