package controller

import (
	"encoding/json"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/common/lib"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dao"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/dto"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/middleware"
	"github.com/WHITE-ILMARE/go_gateway/backend/go_gateway_demo/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	adminController := &AdminController{}
	group.GET("/admin_info", adminController.AdminInfo)
	group.POST("/change_pwd", adminController.ChangePwd)
}

// ChangePwd go doc
// @Summary      修改密码
// @Description  修改密码
// @Tags         管理员接口
// @ID           /admin/change_pwd
// @Produce      json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router       /admin/change_pwd [post]
func (adminController *AdminController) ChangePwd(ctx *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	// 1. session读取用户信息到结构体 sessInfo
	// 2. sessInfo.ID 读取数据库信息 adminInfo
	// 3. saltPassword = sha256(params.password + adminInfo.salt)
	// 4. 数据库保存saltPassword
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	// 从数据库中读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(tx, (&dao.Admin{UserName: adminSessionInfo.UserName}))
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	if err := adminInfo.Save(ctx, tx); err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}
	middleware.ResponseSuccess(ctx, "")
}

// AdminInfo go doc
// @Summary      管理员信息获取
// @Description  管理员信息获取
// @Tags         管理员接口
// @ID           /admin/admin_info
// @Produce      json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router       /admin/admin_info [get]
func (adminController *AdminController) AdminInfo(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	//sessInfoStr := sessInfo.(string)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}
	// 1 读取sessionKey对应json，转换为结构体
	// 2 取出数据然后封装
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://avatars.githubusercontent.com/u/24714274?s=40&v=4",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(ctx, out)
}
