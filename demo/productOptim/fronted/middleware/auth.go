package middleware

import "github.com/kataras/iris/v12"

// AuthConProduct 作为中间件，简单验证用户是否已经登录
func AuthConProduct(ctx iris.Context) {

	uid := ctx.GetCookie("uid")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录!")
		ctx.Redirect("/user/login")	// 重定向到登录界面
		return
	}
	ctx.Application().Logger().Debug("已经登陆")
	ctx.Next()
}
