package dao

import (
	"gorm.io/gorm"
)

type TcpRule struct {
	ID        int `json:"id" gorm:"primary_key"`
	ServiceID int `json:"service_id" gorm:"column:service_id" description:"服务ID"`
	Port      int `json:"port" gorm:"column:port" description:"端口"`
}

func (s *TcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}

func (t *TcpRule) Find(tx *gorm.DB, search *TcpRule) (*TcpRule, error) {
	model := &TcpRule{}
	err := tx.Where(search).First(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *TcpRule) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}
