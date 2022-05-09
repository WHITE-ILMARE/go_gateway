package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	// 基础知识补充
	// go中，一个汉字占3个字节
	testByteArray := []byte("中")
	fmt.Printf("%v\n", testByteArray)
	// len()参数为字符串时返回的也是字节数而不是字符数
	fmt.Printf("len('中')=%d\n", len("中"))
	// 基础知识补充结束

	rs1 := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}
	rs2 := "http://127.0.0.1:2004/base"
	url2, err2 := url.Parse(rs2)
	if err2 != nil {
		log.Println(err2)
	}
	urls := []*url.URL{url1, url2}
	proxy := NewMultipleHostsReverseProxy(urls)
	log.Println("Starting http server at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
