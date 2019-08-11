package cluster

type Node interface {
	Members() []string
	IsProcess(key string) (string, bool)
	HashList() []string
}

var node Node

func New(addr string, cluster []string) {
	var err error
	node, err = NewGossIP(addr, cluster)
	if err != nil {
		panic(err)
	}
}

func Members() []string {
	return node.Members()
}

func IsProcess(key string) (string, bool) {
	return node.IsProcess(key)
}

func HashList() []string {
	return node.HashList()
}
