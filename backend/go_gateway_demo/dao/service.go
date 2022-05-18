package dao

type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"`
	HTTPRule      *HttpRule      `json:"http" description:"http_url"`
	TCPRule       *TcpRule       `json:"tcp" description:"tcp_url"`
	GRPCRule      *GrpcRule      `json:"grpc" description:"grpc"`
	LoadBalance   *LoadBalance   `json:"load_balance" description:"load_balance"`
	AccessControl *AccessControl `json:"access_control" description:"access_control"`
}
