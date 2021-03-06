package dao

import (
	"gorm.io/gorm"
)

type AccessControl struct {
	ID                int    `json:"id" gorm:"primary_key"`
	ServiceID         int    `json:"service_id" gorm:"column:service_id" description:"服务ID"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限，1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流"`
}

func (s *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (t *AccessControl) Find(tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	model := &AccessControl{}
	err := tx.Where(search).Find(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *AccessControl) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}
