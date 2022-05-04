package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	proxy_addr = "http://127.0.0.1:2003"
	port       = "2002"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// step 1: 解析代理地址，并更改请求体的协议和主机
	proxy, err := url.Parse(proxy_addr)
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host

	// step 2: 请求下游
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	if err != nil {
		log.Print(err)
		return
	}

	// step 3: 把下游请求内容返回给上游
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()
	/**
	bufio包提供了有缓冲的I/O，包装一个io.Reader或io.Writer接口对象
	通过缓冲来提高文件读写的效率，io操作本身效率不低，低的是频繁访问磁盘文件，所以bufio提供缓冲区，
	分配一块内存，读写都在缓冲区中，如不命中再读写磁盘文件。例如可以将多次写入操作存储到缓冲区中，最后一次性写入磁盘
	*/
	// resp.Body是ReadCloser结构体实例，实现了Read方法，所以此处用bufio.NewReader创建一个新Reader
	bufio.NewReader(resp.Body).WriteTo(w)
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start serving on port " + port)
	// 在主协程中调用http.ListenAndServe会阻塞主协程，所以主协程不会退出
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
