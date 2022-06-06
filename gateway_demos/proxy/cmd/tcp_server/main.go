package main

// 启动一个TCP服务器
import (
	"context"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/tcp_proxy"
	"log"
	"net"
)

var (
	addr = ":7002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler from 7002\n"))
}

func main() {
	log.Println("Starting tcpserver at " + addr)
	tcpServ := tcp_proxy.TcpServer{
		Addr:    addr,
		Handler: &tcpHandler{},
	}
	fmt.Println("Starting tcp_server at " + addr)
	tcpServ.ListenAndServe()
}
