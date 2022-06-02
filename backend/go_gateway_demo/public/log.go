package public

import (
	"context"
	lib2 "github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/common/lib"
	"github.com/gin-gonic/gin"
)

//错误日志
func ContextWarning(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*lib2.TraceContext)
	if !ok {
		traceContext = lib2.NewTrace()
	}
	lib2.Log.TagWarn(traceContext, dltag, m)
}

//错误日志
func ContextError(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*lib2.TraceContext)
	if !ok {
		traceContext = lib2.NewTrace()
	}
	lib2.Log.TagError(traceContext, dltag, m)
}

//普通日志
func ContextNotice(c context.Context, dltag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*lib2.TraceContext)
	if !ok {
		traceContext = lib2.NewTrace()
	}
	lib2.Log.TagInfo(traceContext, dltag, m)
}

//错误日志
func ComLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	lib2.Log.TagError(traceContext, dltag, m)
}

//普通日志
func ComLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	lib2.Log.TagInfo(traceContext, dltag, m)
}

// 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *lib2.TraceContext {
	// 防御
	if c == nil {
		return lib2.NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*lib2.TraceContext); ok {
			return tc
		}
	}
	return lib2.NewTrace()
}

// 从Context中获取数据
func GetTraceContext(c context.Context) *lib2.TraceContext {
	if c == nil {
		return lib2.NewTrace()
	}
	traceContext := c.Value("trace")
	if tc, ok := traceContext.(*lib2.TraceContext); ok {
		return tc
	}
	return lib2.NewTrace()
}
