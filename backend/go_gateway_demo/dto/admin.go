package dto

import (
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/public"
	"github.com/gin-gonic/gin"
	"time"
)

// 用户信息基于用户Session，但比Session中存储的id、Username和LoginTime多一些属性
type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
