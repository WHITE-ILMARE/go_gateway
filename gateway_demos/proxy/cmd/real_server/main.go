package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rs1 := &RealServer{Addr: "127.0.0.1:2003"}
	rs2 := &RealServer{Addr: "127.0.0.1:2004"}
	// 两个run方法调用http.ListenAndServe()，跑在协程中，不会阻塞主协程
	// 所以需要在主协程中显式声明信号量，等待信号量通知再退出
	rs1.run()
	rs2.run()

	// 监听关闭信号
	quit := make(chan os.Signal)
	// os/signal.Notify函数将监听参数指定的信号集，若有指定信号，就转发到quit中
	// SIGINT: 值为2， 动作为Term，是用户发送INTR字符（Ctrl+C）触发
	// SIGTERM: 值为15，动作为Term，程序结束（可以被捕获、阻塞或忽略）时触发，是kill pid默认发送的信号
	// SIGTERM比kill -9 pid更优雅，因为kill -9 pid向进程发送SIGKILL信号，既不能被应用程序捕获，也不能被阻塞或忽略，立即结束指定进程，在应用程序无感知的情况下被"暴力"终止
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 主协程执行完毕，两个子协程也会退出
	<-quit
}

type RealServer struct {
	Addr string // server要保存自己的addr，供后续使用
}

func (r *RealServer) run() {
	// 用log不用fmt的原因是：log添加了输出时间，线程安全，方便转存日志信息形成日志文件
	log.Println("Starting http server at " + r.Addr)
	// 本质上，ServeMux只是一个路由管理器，它本身也实现了Handler接口的ServeHTTP方法
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloHandler)
	mux.HandleFunc("/base/error", r.ErrorHandler)
	mux.HandleFunc("/test_http_string/aaa", r.TimeoutHandler)
	// 这儿实例化了一个http.Server结构体，写好了各种配置
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	go func() {
		// 这儿就不用指定地址、路由mux了
		log.Fatal(server.ListenAndServe())
	}()
}

// HelloHandler 打印请求的地址
func (r *RealServer) HelloHandler(w http.ResponseWriter, req *http.Request) {
	// 若请求地址为127.0.0.1/abc?q=1，r.Addr为127.0.0.1, Path为/abc
	upath := fmt.Sprintf("Hello from http://%s%s\n", r.Addr, req.URL.Path)

	// 细节是RemoteAddr是在req中内置的字段，而XFF和X-Real-IP是自定义的，只能从Header里取
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-IP=%v\n", req.RemoteAddr, req.Header.Get("X-Forwarded-For"), req.Header.Get("X-Real-IP"))
	headers := fmt.Sprintf("headers = %v\n", req.Header)
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
	io.WriteString(w, headers)
}

func (r *RealServer) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (r *RealServer) TimeoutHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(6 * time.Second)
	msg := "timeout handler"
	w.WriteHeader(200)
	io.WriteString(w, msg)
}
