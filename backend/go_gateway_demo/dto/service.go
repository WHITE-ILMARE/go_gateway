package dto

import (
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/public"
	"github.com/gin-gonic/gin"
)

type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词，搜索用"  validate:""`                            // 关键词
	PageNo   int    `json:"page_no" form:"page_no" comment:"页数" example:"1" validate:""`                // 页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate:"required"` // 每页条数
}

type ServiceListItemOutput struct {
	ID          int    `json:"id" form:"id"`                     // ID
	ServiceName string `json:"service_name" form:"service_name"` // 服务名称
	ServiceDesc string `json:"service_desc" form:"service_desc"` // 服务描述
	LoadType    int    `json:"load_type" form:"load_type"`       // 类型
	ServiceAddr string `json:"service_addr" form:"service_addr"` // service_addr
	Qps         int64  `json:"qps" form:"qps"`                   // qps
	Qpd         int64  `json:"qpd" form:"qpd"`                   // qpd
	TotalNode   int    `json:"total_node" form:"total_node"`     // 节点数
}

type ServiceListOutput struct {
	Total int64                   `json:"total" form:"total" comment:"总数"` // 总数
	List  []ServiceListItemOutput `json:"list" form:"list" comment:"列表"`   // 列表
}

type ServiceDeleteInput struct {
	ID int `json:"id" form:"id" comment:"服务ID" example:"56" validate:"required"` // ID
}

type ServiceAddHTTPInput struct {
	ServiceName            string `json:"service_name" form:"service_name" comment:"服务名"  validate:"required,valid_service_name"`            // 服务名
	ServiceDesc            string `json:"service_desc" form:"service_desc" comment:"服务描述"  validate:"required,max=255,min=1"`                // 服务描述
	RuleType               int    `json:"rule_type" form:"rule_type" comment:"接入类型"  validate:"max=1,min=0"`                                 // 接入类型
	Rule                   string `json:"rule" form:"rule" comment:"接入路径，域名或前缀"  validate:"required,valid_rule"`                             // 接入路径，域名或前缀
	NeedHttps              int    `json:"need_https" form:"need_https" comment:"支持https"  validate:""`                                       // 支持https
	NeedStripUri           int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri"  validate:""`                           // 启用strip_uri
	NeedWebsocket          int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket"  validate:"max=1,min=0"`              // 是否支持websocket
	UrlRewrite             string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能"  validate:"valid_url_rewrite"`                    // url重写功能
	HeaderTransfor         string `json:"header_transfor" form:"header_transfor" comment:"header转换支持"  validate:"valid_header_transfor"`     // header转换支持
	OpenAuth               int    `json:"open_auth" form:"open_auth" comment:"是否开启权限"  validate:"max=1,min=0"`                               // 是否开启权限
	BlackList              string `json:"black_list" form:"black_list" comment:"黑名单ip"  validate:""`                                         // 黑名单ip
	WhiteList              string `json:"white_list" form:"white_list" comment:"白名单ip"  validate:""`                                         // 白名单ip
	ClientIPFlowLimit      int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流"  validate:"min=0"`                // 客户端ip限流
	ServiceFlowLimit       int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端ip限流"  validate:"min=0"`                  // 服务端ip限流
	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式"  validate:"max=3,min=0"`                               // 轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"服务ip列表"  validate:"required,valid_iplist"`                         // 服务ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表"  validate:"required,valid_weight_list"`              // 权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时，单位s"  validate:"min=0"`   // 建立连接超时，单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时，单位s"  validate:"min=0"` // 获取header超时，单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"连接最大空闲时间，单位s"  validate:"min=0"`       // 连接最大空闲时间，单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲连接数"  validate:"min=0"`                    // 最大空闲连接数
}

func (param *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

func (param *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

func (param *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
