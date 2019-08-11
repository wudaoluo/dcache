package cluster

import (
	"fmt"
	"testing"
)

func TestConsistent(t *testing.T) {
	c := NewConsistentHash("127.0.0.1", 256)
	c.Set([]string{"1.1.1.1"})
	fmt.Println(c.HashList())
}
