package dao

import (
	"gorm.io/gorm"
)

type HttpRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务ID"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"type=domain表示域名，type=url_prefix表示url前缀"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"type=1 支持https"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1启用"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"启用websocket 1启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能，每行一个"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header_transfor"`
}

func (s *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (t *HttpRule) Find(tx *gorm.DB, search *HttpRule) (*HttpRule, error) {
	model := &HttpRule{}
	err := tx.Where(search).First(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (t *HttpRule) Save(tx *gorm.DB) error {
	return tx.Save(t).Error
}
