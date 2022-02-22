package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/l1nkkk/shopSystem/demo/product/common"
	"github.com/l1nkkk/shopSystem/demo/product/fronted/middleware"
	"time"

	"github.com/l1nkkk/shopSystem/demo/product/fronted/web/controllers"
	"github.com/l1nkkk/shopSystem/demo/product/repositories"
	"github.com/l1nkkk/shopSystem/demo/product/services"
)

func main() {
	// 1.创建iris 实例
	app := iris.New()
	// 2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	// 3.注册模板
	tmplate := iris.HTML("./demo/product/fronted/web/views", ".html").
		Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	// 4.设置静态资源资源
	//app.StaticWeb("/public", "./fronted/web/public")
	app.HandleDir("/public", "./demo/product/fronted/web/public")
	// 访问生成好的html静态文件
	//app.StaticWeb("/html", "./fronted/web/htmlProductShow")
	app.HandleDir("/html", "./demo/product/fronted/web/htmlProductShow")

	// 5.指定异常跳转页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	// 6.连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {

	}

	// l1nkkk: 登录优化
	//7.创建session
	sess := sessions.New(sessions.Config{
		Cookie:"AdminCookie",		// cookie名称
		Expires:600*time.Minute,	// 过期时间
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 8.注册控制器
	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx,sess.Start)
	//userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("orders", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	// 9.run
	app.Run(
		iris.Addr("0.0.0.0:8082"),
		//iris.WithoutVersionChecker,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
