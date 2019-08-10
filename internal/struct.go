package internal



type Services struct {
	TCP bool
	GRPC bool
	ALL bool
	Listen string
}



type Req struct {
	Op byte
	Retain byte
	Key []byte
	Value []byte
}

func (r *Req) IsPut() bool {
	if r.Op == byte(2) {
		return true
	}
	return false
}