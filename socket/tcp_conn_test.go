package socket

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/wudaoluo/dcache/internal"
)

func TestConn(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:7777")

	if err != nil {
		t.Error("tcpServer.Run", "err", err)
		return
	}

	defer ln.Close()

	go func() {
		time.Sleep(1 * time.Second)
		c1, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			t.Error(err)
			return
		}

		cc := NewTcpConn(c1)
		n, _ := cc.WriteMsg(byte(3), []byte("cccsscscc"), []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbb"))

		fmt.Println("wirite ", n)
		c1.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {

			c := NewTcpConn(conn)

			for {

				var req = internal.Req{}
				time.Sleep(1 * time.Second)
				err = c.ReadMsg(&req)
				if err != nil {
					if err == io.EOF {
						fmt.Println("io.EOF")
						return
					}
					//panic(err)
				}
				fmt.Println(req)
			}
		}()

	}
}
