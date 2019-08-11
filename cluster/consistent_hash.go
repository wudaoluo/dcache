package cluster

import (
	"github.com/wudaoluo/golog"
	"stathat.com/c/consistent"
)

type consistenter interface {
	IsProcess(key string) (string, bool)
	Set(addrs []string)
	HashList() []string
}

type consistentHash struct {
	circle *consistent.Consistent
	addr   string
}

func NewConsistentHash(addr string, numOfReplicas int) consistenter {
	c := &consistentHash{addr: addr}
	c.circle = consistent.New()
	c.circle.NumberOfReplicas = numOfReplicas
	return c
}

func (c *consistentHash) IsProcess(key string) (string, bool) {
	serverIP, err := c.circle.Get(key)
	if err != nil {
		//TODO 触发这日志 说明一致性 hash 有问题
		golog.Error("node.IsProcess", "key", key, "err", err)
		return "", false
	}

	return serverIP, serverIP == c.addr
}

func (c *consistentHash) Set(addrs []string) {
	c.circle.Set(addrs)
}

func (c *consistentHash) HashList() []string {
	return c.circle.Members()
}
