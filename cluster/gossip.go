package cluster

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/memberlist"
)

const NUMBER_OF_REPLICAS = 256

func NewGossIP(addr string, cluster []string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if len(cluster) == 0 {
		cluster = []string{addr}
	}
	_, e = l.Join(cluster)
	if e != nil {
		return nil, e
	}

	n := &gossIP{l, NewConsistentHash(addr, NUMBER_OF_REPLICAS)}
	fmt.Println(l.Members())
	go n.watch()
	return n, nil
}

type gossIP struct {
	l *memberlist.Memberlist
	consistenter
}

func (n *gossIP) Members() []string {
	return n.Members()
}

func (n *gossIP) watch() {
	t := time.NewTicker(3 * time.Second)

	for range t.C {
		m := n.l.Members()
		nodes := make([]string, len(m))
		for i, n := range m {
			nodes[i] = n.Name
		}
		n.Set(nodes)
	}

}
