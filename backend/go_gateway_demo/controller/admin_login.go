package controller

import (
	"encoding/json"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dao"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dto"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/middleware"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLogout)
}

// AdminLogin godoc
// @Summary      管理员登录
// @Description  管理员登录
// @Tags         管理员接口
// @ID           /admin_login/login
// @Accept       json
// @Produce      json
// @Param        body  body  dto.AdminLoginInput  true  "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router       /admin_login/login [post]
func (adminlogin *AdminLoginController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	// 1. params.UserName 取得管理员信息admininfo
	// 2. admininfo.salt + params.Password => saltPassword
	// 3. saltPassword == admininfo.password?
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(tx, params)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	// 设置session
	sessionInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	// sessBts: 二进制数组
	sessBts, err := json.Marshal(sessionInfo)
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	session := sessions.Default(ctx)
	session.Set(public.AdminSessionInfoKey, string(sessBts))
	session.Set("login_test_session", "hhh")
	session.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(ctx, out)
}

// AdminLogout godoc
// @Summary      管理员登出
// @Description  管理员登出
// @Tags         管理员接口
// @ID           /admin_login/logout
// @Accept       json
// @Produce      json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /admin_login/logout [get]
func (adminlogin *AdminLoginController) AdminLogout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(ctx, "")
}
