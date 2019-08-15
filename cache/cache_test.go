package cache

import (
	"log"
	"testing"
)

func TestMemoryCache(t *testing.T) {

	value, ok := Get("a")
	if ok {
		log.Println(value)
	} else {
		log.Println("a", "not found")
	}

	Set("a", 1)

	value, ok = Get("a")
	if ok {
		log.Println("found key a", value)
	} else {
		log.Println("a", "not found")
	}

	Del("a")

}
