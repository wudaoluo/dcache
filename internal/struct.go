package internal

type Services struct {
	TCP     bool
	GRPC    bool
	QUIC    bool
	ALL     bool
	TCPMux  bool
	Listen  string
	MaxConn int
}

type Data struct {
	Op     byte
	Retain byte
	Key    []byte
	Value  []byte
}

func (d *Data) IsValue() bool {
	if d.Op == OP_REQ_PUT || d.Op == OP_REQ_GET {
		return true
	}
	return false
}
