package main

// 实现一个简单的正向代理

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	/**
	req.remoteAddr表示发出请求的远程主机的IP地址，代表客户端的IP，是服务端根据客户端的IP指定的；
	本例中就是发起http请求的地址，访问http://tianya.cn后，发现都是如下格式：
	127.0.0.1:56227
	127.0.0.1:56736
	...
	我理解应该是对于一个标签页访问一个页面，浏览器开了很多http连接请求，用了很多端口
	*/
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	// 数据连接池
	transport := http.DefaultTransport
	// step1: 浅拷贝对象，然后再新增属性数据，避免影响
	outReq := new(http.Request)
	*outReq = *req
	// req.RemoteAddr的格式：IP:Port，是http请求发起方的地址
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// XFF格式：ClientIP, proxy1, proxy2
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			//没太懂，prior是数组吗？为何把clientIP续在最后
			// 解释：可以轻易地查看Header是map[string][]string类型，所以其value都是字符串数组
			// 但是Set时传的key、value都是string，查看Set方法源码得知，背后是将字符串外面套了个数组而已
			// 所以此处取prior[0]亦可
			fmt.Printf("prior=%v\n", prior)
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
			// 写响应头字段
			rw.Header().Add(key, v)
		}
	}
	// 写响应状态码
	rw.WriteHeader(res.StatusCode)
	// 写响应体
	io.Copy(rw, res.Body)
	res.Body.Close()
}

func handleHelloWorld(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	fmt.Println("Serve on :8080")
	http.Handle("/", &Pxy{})
	/**
	服务器上，0.0.0.0代表本机上所有IPV4地址，若一个主机有两个IP地址，若主机上有一监听0.0.0.0的服务，那么两个IP都能访问该服务
	路由中，0.0.0.0代表默认路由；
	本例中，启动的服务监听本机上所有请求
	*/
	http.ListenAndServe("0.0.0.0:8080", nil)
}
