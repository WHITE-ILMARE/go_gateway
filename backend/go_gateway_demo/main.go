package main

import (
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/router"
	"github.com/e421083458/golang_common/lib"
	"os"
	"os/signal"
	"syscall"
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

	lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
	defer lib.Destroy()
	router.HttpServerRun()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	router.HttpServerStop()
}
