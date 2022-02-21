package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/l1nkkk/shopSystem/demo/product/datamodels"
	"github.com/l1nkkk/shopSystem/demo/product/services"

	"strconv"
)

// l1nkkk: 复用后台程序的ProductService
type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}


// GetDetail 获取商品信息； GET /product/detail
func (p *ProductController) GetDetail() mvc.View {
	// l1nkkk: 暂时直接写死
	// 1. 获取商品信息
	product, err := p.ProductService.GetProductByID(2)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	// 2. 模板渲染返回
	return mvc.View{
		Layout: "shared/productLayout.html",	// 不使用main中定义的默认layout
		Name:   "product/view.html",
		Data: iris.Map{		// 传给 template 的数据
			"product": product,
		},
	}
}


//
func (p *ProductController) GetOrder() mvc.View {
	// 1. 通过ID从数据库获取商品信息
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	var orderID int64
	showMessage := "抢购失败！"

	// l1nkkk: 注意这里只是简单的逻辑
	// 2. 判断商品数量是否满足需求
	if product.ProductNum > 0 {
		// 扣除商品数量
		product.ProductNum -= 1
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		// 创建订单
		userID, err := strconv.Atoi(userString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		order := &datamodels.Order{
			UserId:      int64(userID),
			ProductId:   int64(productID),
			OrderStatus: datamodels.OrderSuccess,
		}
		orderID, err = p.OrderService.InsertOrder(order)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}

}
