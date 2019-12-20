package quic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/lucas-clemente/quic-go"
	"github.com/wudaoluo/dcache/internal"
	"github.com/wudaoluo/dcache/service"
	"github.com/wudaoluo/dcache/socket"
	"github.com/wudaoluo/golog"
	"io"
	"math/big"
)

type quicServer struct {
	listen  string
	maxConn int
}

func NewQuicServer(listen string, maxConn int) service.Service {
	return &quicServer{
		listen:  listen,
		maxConn: maxConn,
	}
}

func (q *quicServer) Run() {
	ln, err := quic.ListenAddr(q.listen, generateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	golog.Info("quicServer.Run", "quic server listening addr:", q.listen, "alert", "quic还处于测试阶段,谨慎使用")
	ctx := context.Background()
	for {
		session, err := ln.Accept(ctx)
		if err != nil {
			golog.Error("ln.Accept(ctx)", "err", err)
			return
		}

		go q.handler(ctx, session)
	}

}

func (q *quicServer) handler(ctx context.Context, session quic.Session) {
	defer session.Close()
	for {
		stream, err := session.AcceptStream(ctx)
		if err != nil {
			golog.Error("handler session.AcceptStream", "err", err)
			return
		}

		go q.streamHandler(stream)
	}
}

func (q *quicServer) streamHandler(stream quic.Stream) {
	c := socket.NewQuicStream(stream)
	defer c.Close()
	var err error
	for {
		var data = &internal.Data{}
		err = c.ReadMsg(data)
		if err == io.EOF {
			golog.Info("quicServer.handler", "clientIP", c.RemoteIP(), "err", "io.EOF")
			return
		}
		if err != nil {
			golog.Error("quicServer.handler", "clientIP", c.RemoteIP(), "err", err)
			return
		}

		service.Operate(data)

		_, err = c.WriteMsg(data)
		if err != nil {
			golog.Error("tcpServer.handler", "err", err)
			return
		}
	}

}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
