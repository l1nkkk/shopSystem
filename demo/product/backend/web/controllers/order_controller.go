package controllers

import (

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/l1nkkk/shopSystem/demo/product/services"
)

type OrderController struct { //
	Ctx iris.Context
	OrderService services.IOrderService
}

// Get 获取订单信息;GET /order
func (o *OrderController) Get() mvc.View {
	// 1. 获得所有订单信息
	orderArray,err:=o.OrderService.GetAllOrderInfo()
	if err !=nil {
		o.Ctx.Application().Logger().Debug("查询订单信息失败")
	}

	// 2. 渲染
	return mvc.View{
		Name:"order/view.html",
		Data:iris.Map{
			"order":orderArray,	// 注意模板里使用了range来进行渲染
		},
	}
	
}

