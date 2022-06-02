package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// 非函数内部也可以声明变量并赋值
// 细节是这个地方不能加协议http://，因为addr是传给ListenAndServe函数的参数，指定的是TCP地址
// in the form "host:port". If empty, ":http" (port 80) is used.
var addr = "127.0.0.1:2002"

func main() {
	// 细节是这个地方必须得加前缀http://，不然url.Parse解析不了
	// The url may be relative (a path, without a host) or absolute (starting with a scheme).
	rs1 := "http://127.0.0.1:2003/base"
	url1, err1 := url.Parse(rs1)
	if err1 != nil {
		log.Println(err1)
	}
	proxy := NewSingleHostReverseProxy(url1)
	log.Println("starting http server at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}

// NewSingleHostReverseProxy 新建一个proxy
// 如果target路径是http://127.0.0.1:2003/base
// req的路径如果是http://127.0.0.1:2002/dir
// 则经director修改后的req路径是http://127.0.0.1:2003/base/dir
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	// http://127.0.0.1:2002/dir?name=123
	// RawQuery: name = 123
	// Scheme: http
	// Host: 127.0.0.1
	fmt.Printf("target.RawQuery=%v\n", target.RawQuery)
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		// joinURLPath将参数url的path拼接起来并保证连接处没有多余的‘/’
		// target.Path: /base
		// req.URL.Path: /dir
		// 拼接后：/base/dir
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	modifyfunc := func(res *http.Response) error {
		if res.StatusCode != 200 {
			oldPayload, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			newPayload := []byte("Something unexpected happened, raw content: " + string(oldPayload))
			// res.Body一定要是个ReadCloser，用ioutil.NopCloser构造
			res.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			// 记得要修改Content-length
			res.ContentLength = int64(len(newPayload))
			// 还要同步给Header字段
			res.Header.Set("Content-Length", fmt.Sprint(len(newPayload)))
		}
		return nil
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyfunc}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
