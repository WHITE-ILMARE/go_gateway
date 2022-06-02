package middleware

import (
	"fmt"
	"golang.org/x/time/rate"
)

// RateLimiter 若没有通过限流限制，就终止当前中间件处理过程，直接返回
func RateLimiter() func(c *SliceRouterContext) {
	l := rate.NewLimiter(1, 2)
	return func(c *SliceRouterContext) {
		if !l.Allow() {
			c.Rw.Write([]byte(fmt.Sprintf("rate limit:%v, %v", l.Limit(), l.Burst())))
			c.Abort()
			return
		}
		c.Next()
	}
}
