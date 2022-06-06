package main

import (
	"flag"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/common/lib"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dao"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/http_proxy_router"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/router"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/tcp_proxy_router"
	"os"
	"os/signal"
	"syscall"
)

// 定义从启动命令中读取的变量
var (
	// endpoint 标识是dashboard还是代理服务器
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	// config 标识配置文件夹
	config = flag.String("config", "", "input config file like ./conf/dev/")
)

// @title        网关API
// @version      1.0
// @description  网关API接口文档

// @contact.name   WHITE-ILMARE
// @contact.url    https://github.com/WHITE-ILMARE
// @contact.email  2480800244@qq.com

// @license.name  Apach 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8880
// @BasePath  /

func main() {

	flag.Parse()
	if *endpoint == "" || *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" { // 后台管理服务
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else { // 代理服务器服务
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()
		// 将数据库中的服务信息和APP信息加载到内存中，方便接口调用
		dao.ServiceManagerHandler.LoadOnce()
		dao.AppManagerHandler.LoadOnce()
		// 使用另一套路由服务
		// http代理启动
		go func() {
			http_proxy_router.HttpServerRun()
		}()
		go func() {
			http_proxy_router.HttpsServerRun()
		}()
		// tcp代理启动
		go func() {
			tcp_proxy_router.TcpServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		tcp_proxy_router.TcpServerStop()
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}

}
