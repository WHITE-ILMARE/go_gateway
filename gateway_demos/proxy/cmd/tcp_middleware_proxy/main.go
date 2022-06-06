package main

// 本实例演示使用TCP中间件
import (
	"context"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/load_balance"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/proxy"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/public"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/tcp_middleware"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/tcp_proxy"
	"net"
	"time"
)

var (
	addr = ":2002"
)

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("tcpHandler"))
}

func main() {
	//基于 thrift 代理测试
	rb := load_balance.LoadBanlanceFactory(load_balance.LbWeightRoundRobin)
	rb.Add("127.0.0.1:6001", "40")

	//构建路由及设置中间件
	counter, _ := public.NewFlowCountService("local_app", time.Second)
	router := tcp_middleware.NewTcpSliceRouter()
	router.Group("/").Use(
		tcp_middleware.IpWhiteListMiddleWare(),
		tcp_middleware.FlowCountMiddleWare(counter))

	//构建回调handler
	routerHandler := tcp_middleware.NewTcpSliceRouterHandler(
		func(c *tcp_middleware.TcpSliceRouterContext) tcp_proxy.TCPHandler {
			return proxy.NewTcpLoadBalanceReverseProxy(c, rb)
		}, router)

	//启动服务
	tcpServ := tcp_proxy.TcpServer{Addr: addr, Handler: routerHandler}
	fmt.Println("Starting tcp_proxy at " + addr)
	tcpServ.ListenAndServe()
}
