package cache

type cache interface {
	Get(key interface{}) (interface{},bool)
	Set(key,value interface{})
	Del(key interface{})
}

var _cache cache

func init() {
	_cache = NewMemory()
}

func SetCache(c cache) {
	_cache = c
}

func Get(key interface{}) (interface{},bool) {
	return _cache.Get(key)
}

func Set(key,value interface{})  {
	_cache.Set(key,value)
}

func Del(key interface{}) {
	_cache.Del(key)
}