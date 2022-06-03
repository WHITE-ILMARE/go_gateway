package http_proxy_middleware

import (
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dao"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/middleware"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/reverse_proxy"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//匹配接入方式 基于请求信息
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)
		// 根据服务信息（负载均衡类型、ip列表和权重列表等）构造负载均衡器
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		trans, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		// middleware.ResponseSuccess(c,"ok")
		// return
		// 创建 一个httputil.reverseProxy实例
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		// 作为中间件的最后一层，不必再传递，直接使用ServeHTTP处理真实的HTTP请求
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}
