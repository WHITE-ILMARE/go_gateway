package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
)

type HandlerFunc func(*SliceRouterContext)

const abortIndex int8 = math.MaxInt8 / 2 // 最多63个中间件

type SliceRouter struct {
	groups []*SliceGroup
}

func NewSliceRouter() *SliceRouter {
	return &SliceRouter{}
}

// Group 为传入的路径创建具体路由，并绑定其所属的SliceRouter，为use方法做铺垫
func (g *SliceRouter) Group(path string) *SliceGroup {
	return &SliceGroup{
		SliceRouter: g,
		path:        path,
	}
}

// SliceGroup 应该是指明路由和该路由上的中间件
// 一个SliceGroup代表一个路由
// 为什么要有*SliceRouter?访问和自己同组的路由吗？
type SliceGroup struct {
	*SliceRouter
	path     string
	handlers []HandlerFunc
}

// Use 把自己插入调用者SliceGroup的handlers数组中，调用者再检查自己是否还未插入SliceRouter的groups中
// 疑问：这样每次use都检查性能太低了吧
func (g *SliceGroup) Use(middlewares ...HandlerFunc) *SliceGroup {
	g.handlers = append(g.handlers, middlewares...)
	existsFlag := false
	// SliceRouter->当前路由器，SliceGroup->一个具体路由，可能会没有自己吗？
	for _, oldGroup := range g.SliceRouter.groups {
		if oldGroup == g {
			existsFlag = true
		}
	}
	if !existsFlag {
		g.SliceRouter.groups = append(g.SliceRouter.groups, g)
	}
	return g
}

// SliceRouterContext 存了路由器可能用到的所有外部信息
type SliceRouterContext struct {
	// http包相关内容
	Rw  http.ResponseWriter
	Req *http.Request
	Ctx context.Context // 构造时使用*http.Request.context()
	// SliceRouterContext 包含了一个SliceGroup具体路由
	*SliceGroup
	index int8
}

// newSliceRouterContext 创建路由上下文
func newSliceRouterContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}
	// 找最长url前缀匹配
	matchUrlLen := 0
	for _, group := range r.groups {
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				*newSliceGroup = *group // 浅拷贝指针
			}
		}
	}
	c := &SliceRouterContext{
		Rw:         rw,
		Req:        req,
		SliceGroup: newSliceGroup,
		Ctx:        req.Context(),
	}
	return c
}

func (c *SliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

// Set 直接设置整个Ctx，而不是往Ctx中新增k-v，看来一个Context只能存一对k-v
// 按照WithValue的用法，key需要是一个自定义的类型，不能是built-in类型
func (c *SliceRouterContext) Set(key, value interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, value)
}

// Next 从当前index的下一个位置执行至结尾？
func (c *SliceRouterContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 跳出中间件方法
func (c *SliceRouterContext) Abort() {
	c.index = abortIndex
}

func (c *SliceRouterContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *SliceRouterContext) Reset() {
	c.index = -1
}

// SliceRouterHandler 自身就是一个http.Handler
type SliceRouterHandler struct {
	// http.Handler 是接口，只有一个ServeHTTP(ResponseWriter, *Request)方法
	// 所以coreFunc就是接受一个路由器上下文，返回一个能处理http请求的接口
	coreFunc func(routerContext *SliceRouterContext) http.Handler
	router   *SliceRouter
}

func (w *SliceRouterHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := newSliceRouterContext(rw, req, w.router)
	if w.coreFunc != nil {
		c.handlers = append(c.handlers, func(c *SliceRouterContext) {
			w.coreFunc(c).ServeHTTP(rw, req)
		})
	}
	c.Reset()
	c.Next()
}

func NewSliceRouterHandler(coreFunc func(*SliceRouterContext) http.Handler, router *SliceRouter) *SliceRouterHandler {
	return &SliceRouterHandler{
		coreFunc: coreFunc,
		router:   router,
	}
}
