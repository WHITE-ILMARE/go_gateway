package dao

import (
	"gorm.io/gorm"
)

type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务ID"`
	Port           int    `json:"port" gorm:"column:port" description:"端口"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add),删除（del），修改（edit）格式：add headername headervalue"`
}

func (s *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (t *GrpcRule) Find(tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	model := &GrpcRule{}
	err := tx.Where(search).Find(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *GrpcRule) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}
