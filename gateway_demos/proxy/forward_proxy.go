package main

// 实现一个简单的正向代理

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	transport := http.DefaultTransport
	// step1: 浅拷贝对象，然后再新增属性数据
	outReq := new(http.Request)
	*outReq = *req
	// req.RemoteAddr的格式：IP:Port，是http请求发起方的地址
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// XFF格式：ClientIP, proxy1, proxy2
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			//没太懂，prior是数组吗？为何把clientIP续在最后
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	// step2:请求下游
	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step3:把下游响应返回给上游
	for key, value := range res.Header {
		for _, v := range value {
			// 把下游响应的Header头内容加到给上游的响应中
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func handleHelloWorld(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	http.HandleFunc("/", handleHelloWorld)
	log.Fatal(http.ListenAndServe(":4321", nil))
}
