package cache

import "sync"

type memoryCache struct {
	sync.Map
}

func NewMemory() cache {
	return &memoryCache{

	}
}

func (mc *memoryCache) Get(key interface{}) (interface{},bool) {
	return mc.Load(key)
}


func (mc *memoryCache) Set(key,value interface{}) {
	mc.Store(key,value)
}


func (mc *memoryCache) Del(key interface{}) {
	mc.Delete(key)
}

