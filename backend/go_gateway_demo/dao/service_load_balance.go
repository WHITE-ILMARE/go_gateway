package dao

import (
	"gorm.io/gorm"
	"strings"
)

type LoadBalance struct {
	ID                     int64  `json:"id" gorm:"primary_key"`
	ServiceID              int64  `json:"service_id" gorm:"column:service_id" description:"服务ID"`
	CheckMethod            int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功"`
	CheckTimeout           int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间"`
	CheckInterval          int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔，单位s"`
	RoundType              int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList                 string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList             string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList             string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"upstream_connect_timeout"`
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"upstream_header_timeout"`
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"upstream_idle_timeout"`
	UpstreamMaxIdle        int    `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"upstream_max_idle"`
}

func (s *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) Find(tx *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	model := &LoadBalance{}
	err := tx.Where(search).Find(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *LoadBalance) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}

func (t *LoadBalance) GetIPListByModel() []string {
	return strings.Split(t.IpList, ",")
}
