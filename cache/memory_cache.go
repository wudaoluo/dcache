package cache

import (
	"sync"

	"github.com/wudaoluo/golog"
)

type memoryCache struct {
	sync.Map
}

func NewMemory() cache {
	return &memoryCache{}
}

func (mc *memoryCache) Get(key []byte) ([]byte, bool) {
	buf, state := mc.Load(string(key))

	if !state {
		return nil, false
	}

	b, ok := buf.([]byte)
	if ok {
		return b, state
	}

	golog.Warn("memoryCache.Get 断言", "key", string(key))
	return nil, false

}

func (mc *memoryCache) Set(key, value []byte) {
	mc.Store(string(key), value)
}

func (mc *memoryCache) Del(key []byte) {
	mc.Delete(string(key))
}
