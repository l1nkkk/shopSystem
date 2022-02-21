package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/l1nkkk/shopSystem/demo/productOptim/datamodels"
	"github.com/l1nkkk/shopSystem/demo/productOptim/services"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

// l1nkkk: 复用后台程序的ProductService
type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

// =========l1nkkk: 前端优化关键 ==Start

// 以 templatePath 中的模板文件为模板，渲染出静态页面放在 htmlOutPath 中
var (
	// 生成的Html保存目录
	htmlOutPath = "./demo/productOptim/fronted/web/htmlProductShow/"
	// 静态文件模版目录
	templatePath = "./demo/productOptim/fronted/web/views/template/"
)

// GetGenerateHtml 根据模板和SQL中的数据，生成html静态文件；
// GET product/generate/html?productID=1
func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID,err:=strconv.Atoi(productString)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 1.获取模版
	contenstTmp,err := template.ParseFiles(filepath.Join(templatePath,"product.html"))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 2.获取html生成路径
	fileName:=filepath.Join(htmlOutPath,"htmlProduct.html")

	// 3.获取模版渲染数据（数据库操作）
	product,err:=p.ProductService.GetProductByID(int64(productID))
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	// 4.生成静态文件
	generateStaticHtml(p.Ctx,contenstTmp,fileName,product)
}

// generateStaticHtml 生成html静态文件
func generateStaticHtml(ctx iris.Context,template *template.Template,fileName string,product *datamodels.Product)  {
	// 1.判断静态文件是否存在
	if exist(fileName) {
		// 如果存在则删除？？？，那还有啥用
		err:=os.Remove(fileName)
		if err !=nil {
			ctx.Application().Logger().Error(err)
		}
	}
	// 2.生成静态文件
	file,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,os.ModePerm)
	if err !=nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file,&product)
}

// exist 判断文件是否存在，目的是为了判断静态文件是否已经生成
func exist(fileName string) bool  {
	_,err:=os.Stat(fileName)
	return err==nil || os.IsExist(err)
}
// =========l1nkkk: 前端优化关键 ==End

// GetDetail 获取商品详细信息； GET /product/detail
func (p *ProductController) GetDetail() mvc.View {
	// l1nkkk: 暂时直接写死
	// 1. 获取商品信息
	product, err := p.ProductService.GetProductByID(2)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	// 2. 模板渲染返回
	return mvc.View{
		Layout: "shared/productLayout.html", // 不使用main中定义的默认layout
		Name:   "product/view.html",
		Data: iris.Map{ // 传给 template 的数据
			"product": product,
		},
	}
}


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
