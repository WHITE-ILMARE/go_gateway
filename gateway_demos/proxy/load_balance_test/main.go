package main

import (
	"bytes"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/load_balance"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

var addr = "127.0.0.1:2002"

func NewMultipleHostsReverseProxy(lb load_balance.LoadBalance) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		// 根据整个URL计算的hash，只有URL完全一样才一定会分发给同一个服务器
		// 如果想根据IP，就可以取req.RemoteAddr
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil {
			log.Fatal("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawPath = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	// 更改内容
	modifyFunc := func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			newPayload := []byte("StatusCode != 200!!! and oldPayload is: " + string(oldPayload))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			resp.ContentLength = int64(len(newPayload))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
		}
		return nil
	}

	// 错误回调
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func main() {
	rb := load_balance.LoadBalanceFactory(load_balance.LbConsistentHash)
	if err := rb.Add("http://127.0.0.1:2003/base", "10"); err != nil {
		log.Println(err)
	}
	if err := rb.Add("http://127.0.0.1:2004/base", "20"); err != nil {
		log.Println(err)
	}
	proxy := NewMultipleHostsReverseProxy(rb)
	log.Println("Starting http server at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
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
