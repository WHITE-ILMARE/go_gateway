package main

// 本实例演示TCP代理各种上层协议
import (
	"context"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/load_balance"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/proxy"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/tcp_middleware"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/tcp_proxy"
	"net"
)

var (
	addr = ":2002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler\n"))
}

func main() {
	rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)

	// tcp/thrift代理测试
	// tcp真实服务器地址7002
	//rb.Add("127.0.0.1:7002", "100")
	// thrift服务地址6001
	//rb.Add("127.0.0.1:6001", "100")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//redis服务器测试
	//rb.Add("127.0.0.1:6379", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//http服务器测试:
	//缺点对请求的管控不足,比如我们用来做baidu代理,因为无法更改请求host,所以很轻易把我们拒绝
	//rb.Add("127.0.0.1:2003", "40")
	//rb.Add("127.0.0.1:2004", "40")
	rb.Add("www.baidu.com:80", "40")
	proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy}
	fmt.Println("tcp_proxy start at:" + addr)
	tcpServ.ListenAndServe()

	//websocket服务器测试:缺点对请求的管控不足
	//rb.Add("127.0.0.1:2003", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()

	//http2服务器测试:缺点对请求的管控不足
	//rb.Add("127.0.0.1:3003", "40")
	//proxy := proxy.NewTcpLoadBalanceReverseProxy(&tcp_middleware.TcpSliceRouterContext{}, rb)
	//tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: proxy,}
	//fmt.Println("Starting tcp_proxy at " + addr)
	//tcpServ.ListenAndServe()
}
