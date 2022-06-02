package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	// 下游服务器地址
	rs1 := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}
	// 创建反向代理，将本服务器接收到的请求转发到url1并返回响应(addr->url1)
	// 会将我们请求addr的path也append到url1后面
	proxy := httputil.NewSingleHostReverseProxy(url1)
	log.Println("Starting http server at " + addr)
	// proxy实现了ServeHTTP方法，所以它直接可以作为handler传入ListenAndServe中
	log.Fatal(http.ListenAndServe(addr, proxy))
}
