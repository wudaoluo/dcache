package cache

type cache interface {
	Get(key []byte) ([]byte, bool)
	Set(key, value []byte)
	Del(key []byte)
}

var _cache cache

func init() {
	_cache = NewMemory()
}

func SetCache(c cache) {
	_cache = c
}

func Get(key []byte) ([]byte, bool) {
	return _cache.Get(key)
}

func Set(key, value []byte) {
	_cache.Set(key, value)
}

func Del(key []byte) {
	_cache.Del(key)
}
