package main

import (
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/load_balance"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/middleware"
	proxy2 "github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/proxy"
	"log"
	"net/http"
)

var (
	addr_websocket_proxy_main = "127.0.0.1:2002"
)

func main() {
	rb := load_balance.LoadBalanceFactory(load_balance.LbWeightRoundRobin)
	rb.Add("http://127.0.0.1:2003", "50")
	proxy := proxy2.NewLoadBalanceReverseProxy(&middleware.SliceRouterContext{}, rb)
	log.Println("Starting httpserver at " + addr_websocket_proxy_main)
	log.Fatal(http.ListenAndServe(addr_websocket_proxy_main, proxy))
}
