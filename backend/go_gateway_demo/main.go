package main

import (
	"flag"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/lib"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/router"
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
		router.HttpServerRun()

		fmt.Println("启动代理服务了！ ")

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	}

}
