package main

import (
	"context"

	"github.com/l1nkkk/shopSystem/demo/productOptim/backend/web/controllers"
	"github.com/l1nkkk/shopSystem/demo/productOptim/common"
	"github.com/l1nkkk/shopSystem/demo/productOptim/repositories"
	"github.com/l1nkkk/shopSystem/demo/productOptim/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	// 1.创建iris 实例
	app := iris.New()

	// 2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	// 3.注册模板（模板的目录，模板的文件后缀）
	// Layout指定布局文件
	tmplate := iris.HTML("./demo/productOptim/backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	// 4-1.设置模板静态资源目录
	//app.StaticWeb("/assets", "./backend/web/assets") // 旧方法，已经弃用
	app.HandleDir("/assets", "./demo/productOptim/backend/web/assets")
	// 4-2.指定异常跳转页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		// 要传的错误信息
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		// 设置layout
		ctx.ViewLayout("")
		// 渲染的视图
		ctx.View("shared/error.html")
	})

	// l1nkkk: 用于构造 productRepository，也可以不用连接数据库这一步，
	// productRepository 调用 Conn的时候，会自动建立好连接。
	// 4-3.连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// l1nkkk: productSerivce 和 ProductController 怎么建立联系的，这里有点抽象
	// 5-1. 注册控制器——product
	// model
	productRepository := repositories.NewProductManager("product", db)
	productSerivce := services.NewProductService(productRepository)
	productParty := app.Party("/product")              // url 服务前缀
	product := mvc.New(productParty)                   // mvc 对象
	product.Register(ctx, productSerivce)              // 注册model（service）
	product.Handle(new(controllers.ProductController)) // 注册controller

	// 5-2. 注册控制器——Order
	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order") // 将该前缀的url都定位到该controller
	order := mvc.New(orderParty)
	order.Register(ctx, orderService) // 对应了 OrderController 的成员
	order.Handle(new(controllers.OrderController))

	// 6.启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		//iris.WithoutVersionChecker,                    // 启动时是否启动iris版本, 已弃用
		iris.WithoutServerError(iris.ErrServerClosed), // 忽略iris框架的错误
		iris.WithOptimizations,                        // what？
	)

}
