package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"

	"flag"
	"time"

	"sync"

	"bytes"

	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/dcache/socket"
)

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

func init() {
	//flag.StringVar(&typ, "type", "redis", "cache server type")
	flag.StringVar(&server, "h", "localhost:7777", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.StringVar(&operation, "t", "put", "test set, could be get/put/del")
	flag.IntVar(&keyspacelen, "r", 10000000000, "keyspacelen, use random keys from 0 to keyspacelen-1")
	//flag.IntVar(&pipelen, "P", 1, "pipeline length")
	flag.Parse()
	//fmt.Println("type is", typ)
	fmt.Println("server is", server)
	fmt.Println("total", total, "requests")
	fmt.Println("data size is", valueSize)
	fmt.Println("we have", threads, "connections")
	fmt.Println("operation is", operation)
	fmt.Println("keyspacelen is", keyspacelen)
	//fmt.Println("pipeline length is", pipelen)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	ch := make(chan int, 100)
	go func() {
		for i := 0; i < total; i++ {
			ch <- i
		}
		close(ch)
	}()

	var wg = sync.WaitGroup{}
	start := time.Now()
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go process(ch, &wg)
	}

	wg.Wait()
	d := time.Now().Sub(start)

	fmt.Printf("%f seconds total\n", d.Seconds())
	fmt.Printf("rps is %f\n", float64(total)/float64(d.Seconds()))
	fmt.Printf("throughput is %f MB/s\n", +float64(total*valueSize)/1e6/d.Seconds())
}

func process(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", server)
	if err != nil {
		panic(err)
	}

	c := socket.NewTcpConn(conn)

	var op byte
	switch operation {
	case "put":
		op = internal.OP_REQ_PUT
	case "get":
		op = internal.OP_REQ_GET
	case "del":

		op = internal.OP_REQ_DEL

	}

	var data = &internal.Data{}
	value := bytes.Repeat([]byte("a"), valueSize)
	for {
		_, ok := <-ch
		if !ok {
			c.Close()
			//golog.Warn("ch close")
			return
		}

		data.Op = op
		data.Key = []byte(strconv.Itoa(rand.Intn(keyspacelen)))
		data.Value = value
		_, err = c.WriteMsg(data)

		data.Value = nil
		err = c.ReadMsg(data)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}
