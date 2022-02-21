package controllers

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/l1nkkk/shopSystem/demo/productOptim/datamodels"
	"github.com/l1nkkk/shopSystem/demo/productOptim/services"
	"github.com/l1nkkk/shopSystem/demo/productOptim/tool"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

// GetRegister 返回 Register 页面；GET /user/register
func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

// PostRegister 处理用户提交的注册信息表单； POST /user/register
func (c *UserController) PostRegister() {
	// 1. 解析表单，适用于字段较少的情况，对比product的处理
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	// l1nkkk: 可以通过以下扩展来验证表单
	// ozzo-validation；
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}

	// 2. 插入user
	_, err := c.Service.AddUser(user)
	c.Ctx.Application().Logger().Debug(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return // 必须要Redirect后return
	}
	c.Ctx.Redirect("/user/login")
	return
}

// GetLogin 处理登录页面请求；GET /user/login
func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

// PostLogin 处理登录信息表单；POST /user/login
func (c *UserController) PostLogin() mvc.Response {
	//1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	// 2、验证账号密码正确
	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// 3、写入用户ID到cookie中
	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))

	// 重定向到商品信息页面
	return mvc.Response{
		Path: "/product/",
	}

}
