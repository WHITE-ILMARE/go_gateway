package controller

import (
	"errors"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dao"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dto"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/lib"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/middleware"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/public"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type ServiceController struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &ServiceController{}
	group.GET("/service_list", service.ServiceList)
	group.GET("/service_delete", service.ServiceDelete)
	group.GET("/service_detail", service.ServiceDetail)
	group.GET("/service_stat", service.ServiceStat)
	group.POST("/service_add_http", service.ServiceAddHTTP)
	group.POST("/service_update_http", service.ServiceUpdateHTTP)
	group.POST("/service_add_grpc", service.ServiceAddGRPC)
	group.POST("/service_update_grpc", service.ServiceUpdateGRPC)
	group.POST("/service_add_tcp", service.ServiceAddTcp)
	group.POST("/service_update_tcp", service.ServiceUpdateTcp)
}

// ServiceList go doc
// @Summary      服务列表
// @Description  服务列表
// @Tags         服务管理
// @ID           /service/service_list
// @Produce      json
// @Param info query string false "关键词"
// @Param page_size query int true "每页个数"
// @Param page_no query int true "当前页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router       /service/service_list [get]
func (service *ServiceController) ServiceList(ctx *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 从db中分页读取基本信息
	serviceInfo := &dao.ServiceInfo{}
	list, total, err := serviceInfo.PageList(tx, params)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	// 格式化输出信息
	outList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(tx, &listItem)
		if err != nil {
			middleware.ResponseError(ctx, 2003, err)
			return
		}
		// 1. http后缀接入 clusterIP + clusterPort + path
		// 2. http域名接入 domain
		// 3. tcp,grpc接入 clusterIP+servicePort
		serviceAddr := "unknown"
		clusterIP := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}

		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIPListByModel()
		outItem := dto.ServiceListItemOutput{
			ID:          listItem.ID,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qpd:         0,
			Qps:         0,
			TotalNode:   len(ipList),
			LoadType:    listItem.LoadType,
		}
		outList = append(outList, outItem)
	}

	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}
	middleware.ResponseSuccess(ctx, out)
}

// ServiceDelete go doc
// @Summary      服务删除
// @Description  服务删除
// @Tags         服务管理
// @ID           /service/service_delete
// @Produce      json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /service/service_delete [get]
func (service *ServiceController) ServiceDelete(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 读取服务基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	serviceInfo.IsDelete = 1
	if err = serviceInfo.Save(tx); err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	middleware.ResponseSuccess(ctx, "")
}

// ServiceDetail go doc
// @Summary      服务详情
// @Description  服务详情
// @Tags         服务管理
// @ID           /service/service_detail
// @Produce      json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router       /service/service_detail [get]
func (service *ServiceController) ServiceDetail(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 读取服务基本信息
	serviceInfo := &dao.ServiceInfo{ID: params.ID}
	serviceInfo, err = serviceInfo.Find(tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	middleware.ResponseSuccess(ctx, serviceDetail)
}

// ServiceStat go doc
// @Summary      服务统计
// @Description  服务统计
// @Tags         服务管理
// @ID           /service/service_stat
// @Produce      json
// @Param id query string true "服务ID"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router       /service/service_stat [get]
func (service *ServiceController) ServiceStat(ctx *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	//tx, err := lib.GetGormPool("default")
	//if err != nil {
	//	middleware.ResponseError(ctx, 2001, err)
	//	return
	//}
	// 读取服务基本信息
	//serviceInfo := &dao.ServiceInfo{ID: params.ID}
	//serviceInfo, err = serviceInfo.Find(tx, serviceInfo)
	//if err != nil {
	//	middleware.ResponseError(ctx, 2002, err)
	//	return
	//}
	//serviceDetail, err := serviceInfo.ServiceDetail(tx, serviceInfo)
	//if err != nil {
	//	middleware.ResponseError(ctx, 2003, err)
	//	return
	//}
	var todayList []int
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	var yesterdayList []int
	for i := 0; i <= 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}
	middleware.ResponseSuccess(ctx, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}

// ServiceAddHTTP go doc
// @Summary      添加HTTP服务
// @Description  添加HTTP服务
// @Tags         服务管理
// @ID           /service/service_add_http
// @Produce      json
// @Param body body dto.ServiceAddHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /service/service_add_http [post]
func (serviceController *ServiceController) ServiceAddHTTP(ctx *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 跨字段的校验
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(ctx, 2004, errors.New("IP列表与权重列表数量不一致"))
		return
	}
	tx = tx.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	// 按ServiceName查询
	if _, err = serviceInfo.Find(tx, serviceInfo); err == nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2002, errors.New("服务已存在"))
		return
	}
	// 基于数据库的校验
	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Find(tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2003, errors.New("服务接入前缀或域名已存在"))
		return
	}
	// 先向service_info表存，并取得新插入条目的ID
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, err)
		return
	}
	// 此时已经可以拿到ID
	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	// http_rule表创建完毕
	if err := httpRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	// 接下来处理权限控制表
	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2007, err)
		return
	}
	// 接下来插入load-balance表
	loadBalance := &dao.LoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err = loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2008, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
}

// ServiceUpdateHTTP go doc
// @Summary      修改HTTP服务
// @Description  修改HTTP服务
// @Tags         服务管理
// @ID           /service/service_update_http
// @Produce      json
// @Param body body dto.ServiceUpdateHTTPInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /service/service_update_http [post]
func (serviceController *ServiceController) ServiceUpdateHTTP(ctx *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 跨字段的校验
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(ctx, 2002, errors.New("IP列表与权重列表数量不一致"))
		return
	}
	tx = tx.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	// 因为serviceInfo.ServiceDetail直接将传入的serviceInfo作为其Info字段的值了
	// 而我们构造的serviceInfo还不完整，所以需要Find一次，拿到完整字段
	serviceInfo, err = serviceInfo.Find(tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(tx, serviceInfo)
	// 按ServiceName查询
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2004, errors.New("服务不存在"))
		return
	}
	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err = info.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, err)
		return
	}
	// 先更新http_rule表，ServiceID,RuleType和Rule字段不能调整，其余字段可调整
	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	// http_rule表更新完毕
	if err := httpRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, err)
		return
	}
	// 接下来处理权限控制表
	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	// 接下来更新load-balance表
	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err = loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2007, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
}

// ServiceAddGRPC go doc
// @Summary      添加GRPC服务
// @Description  添加GRPC服务
// @Tags         服务管理
// @ID           /service/service_add_grpc
// @Produce      json
// @Param body body dto.ServiceAddGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /service/service_add_grpc [post]
func (serviceController *ServiceController) ServiceAddGRPC(ctx *gin.Context) {
	params := &dto.ServiceAddGrpcInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 验证service_name是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(tx, infoSearch); err == nil {
		middleware.ResponseError(ctx, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}
	// 验证端口是否被占用
	// 验证这个端口与TCP服务是否冲突
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err = tcpRuleSearch.Find(tx, tcpRuleSearch); err == nil {
		middleware.ResponseError(ctx, 2003, errors.New("端口被占用，请重新输入"))
		return
	}
	// 验证这个端口与GRPC服务是否冲突
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err = grpcRuleSearch.Find(tx, grpcRuleSearch); err == nil {
		middleware.ResponseError(ctx, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}
	// 跨字段的校验,ip列表与权重列表匹配
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, errors.New("IP列表与权重列表数量不一致"))
		return
	}
	tx = tx.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	// 接下来插入load-balance表
	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err = loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2007, err)
		return
	}
	// 接下来插入sercice_grpc表
	grpcRule := &dao.GrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeadTransfor,
	}
	if err = grpcRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2008, err)
		return
	}
	// 插入access_control表
	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
	return
}

// ServiceUpdateGRPC go doc
// @Summary      修改GRPC服务
// @Description  修改GRPC服务
// @Tags         服务管理
// @ID           /service/service_update_grpc
// @Produce      json
// @Param body body dto.ServiceUpdateGrpcInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /service/service_update_grpc [post]
func (serviceController *ServiceController) ServiceUpdateGRPC(ctx *gin.Context) {
	params := &dto.ServiceUpdateGrpcInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	// 数据库连接
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// ip列表与权重列表一致
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(ctx, 2002, errors.New("IP列表与权重列表数量不一致"))
		return
	}
	tx = tx.Begin()
	service := &dao.ServiceInfo{ID: params.ID}
	serviceDetail, err := service.ServiceDetail(tx, service)
	if err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2003, errors.New("服务不存在"))
		return
	}
	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err = info.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2004, err)
		return
	}
	// 更新load-balance表
	loadBalance := &dao.LoadBalance{}
	// 有的grpc服务可能没有load-balance部分
	if serviceDetail.LoadBalance != nil {
		loadBalance = serviceDetail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err = loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, err)
		return
	}
	// 处理grpc_rule表
	grpcRule := &dao.GrpcRule{}
	if serviceDetail.GRPCRule != nil {
		grpcRule = serviceDetail.GRPCRule
	}
	grpcRule.ServiceID = info.ID
	grpcRule.HeaderTransfor = params.HeadTransfor
	if err := grpcRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	// 处理权限控制表
	accessControl := &dao.AccessControl{}
	if serviceDetail.AccessControl != nil {
		accessControl = serviceDetail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}

	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
}

// ServiceAddTcp godoc
// @Summary tcp服务添加
// @Description tcp服务添加
// @Tags 服务管理
// @ID /service/service_add_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceAddTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_add_tcp [post]
func (admin *ServiceController) ServiceAddTcp(ctx *gin.Context) {
	params := &dto.ServiceAddTcpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2005, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	//验证 service_name 是否被占用
	infoSearch := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err := infoSearch.Find(tx, infoSearch); err == nil {
		middleware.ResponseError(ctx, 2002, errors.New("服务名被占用，请重新输入"))
		return
	}
	//验证端口是否被占用?
	tcpRuleSearch := &dao.TcpRule{
		Port: params.Port,
	}
	if _, err := tcpRuleSearch.Find(tx, tcpRuleSearch); err == nil {
		middleware.ResponseError(ctx, 2003, errors.New("服务端口被占用，请重新输入"))
		return
	}
	grpcRuleSearch := &dao.GrpcRule{
		Port: params.Port,
	}
	if _, err := grpcRuleSearch.Find(tx, grpcRuleSearch); err == nil {
		middleware.ResponseError(ctx, 2004, errors.New("服务端口被占用，请重新输入"))
		return
	}
	tx = tx.Begin()
	info := &dao.ServiceInfo{
		LoadType:    public.LoadTypeTCP,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := info.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	loadBalance := &dao.LoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err := loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2007, err)
		return
	}
	httpRule := &dao.TcpRule{
		ServiceID: info.ID,
		Port:      params.Port,
	}
	if err := httpRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2008, err)
		return
	}

	accessControl := &dao.AccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIPFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2009, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
	return
}

// ServiceUpdateTcp godoc
// @Summary tcp服务更新
// @Description tcp服务更新
// @Tags 服务管理
// @ID /service/service_update_tcp
// @Accept  json
// @Produce  json
// @Param body body dto.ServiceUpdateTcpInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_update_tcp [post]
func (admin *ServiceController) ServiceUpdateTcp(ctx *gin.Context) {
	params := &dto.ServiceUpdateTcpInput{}
	if err := params.GetValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	//ip与权重数量一致
	if len(strings.Split(params.IpList, ",")) != len(strings.Split(params.WeightList, ",")) {
		middleware.ResponseError(ctx, 2002, errors.New("ip列表与权重设置不匹配"))
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	//tx := lib.GORMDefaultPool.Begin()
	service := &dao.ServiceInfo{
		ID: params.ID,
	}
	detail, err := service.ServiceDetail(tx, service)
	if err != nil {
		middleware.ResponseError(ctx, 2004, err)
		return
	}
	// 更新service_info表，只修改ServiceDesc字段
	info := detail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2005, err)
		return
	}
	// 更新load_balance表
	loadBalance := &dao.LoadBalance{}
	if detail.LoadBalance != nil {
		loadBalance = detail.LoadBalance
	}
	loadBalance.ServiceID = info.ID
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.ForbidList = params.ForbidList
	if err := loadBalance.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2006, err)
		return
	}
	// 更新tcp_rule表
	tcpRule := &dao.TcpRule{}
	if detail.TCPRule != nil {
		tcpRule = detail.TCPRule
	}
	tcpRule.ServiceID = info.ID
	tcpRule.Port = params.Port
	if err := tcpRule.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2007, err)
		return
	}
	// 更新access_control表
	accessControl := &dao.AccessControl{}
	if detail.AccessControl != nil {
		accessControl = detail.AccessControl
	}
	accessControl.ServiceID = info.ID
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.WhiteHostName = params.WhiteHostName
	accessControl.ClientIPFlowLimit = params.ClientIPFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := accessControl.Save(tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(ctx, 2008, err)
		return
	}
	tx.Commit()
	middleware.ResponseSuccess(ctx, "")
	return
}
