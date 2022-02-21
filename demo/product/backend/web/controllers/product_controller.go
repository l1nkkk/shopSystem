package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/l1nkkk/shopSystem/demo/product/common"
	"github.com/l1nkkk/shopSystem/demo/product/datamodels"
	"github.com/l1nkkk/shopSystem/demo/product/services"

	"strconv"
)

type ProductController struct {
	Ctx iris.Context
	ProductService services.IProductService
}


// l1nkkk: 这里定义的方法开头有讲究，需要RestFul支持的操作

// GetAll 获取所有商品；Get product/all
func (p *ProductController) GetAll() mvc.View {
	productArray ,_:=p.ProductService.GetAllProduct()
	return mvc.View{
		Name:"product/view.html",
		Data:iris.Map{
			"productArray":productArray,
		},
	}
}

// PostUpdate 修改商品；Post product/update
func (p *ProductController) PostUpdate ()  {
	product :=&datamodels.Product{}
	// 1. 解析表单
	p.Ctx.Request().ParseForm()

	// 2. 将表单映射到product对象
	// l1nkkk: 这里用到了 imooc标签，通过 tag 将表单映射到product结构体
	dec := common.NewDecoder(&common.DecoderOptions{TagName:"imooc"})
	if err:= dec.Decode(p.Ctx.Request().Form,product);err!=nil {
		// 打印日志
		p.Ctx.Application().Logger().Debug(err)
	}

	// 3. 调用Model（service）进行处理（更新数据）
	err:=p.ProductService.UpdateProduct(product)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 4. 跳转界面
	p.Ctx.Redirect("/product/all")
}





func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name:"product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	// 遇上面的PostUpdate基本一样
	product :=&datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName:"imooc"})
	if err:= dec.Decode(p.Ctx.Request().Form,product);err!=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_,err:=p.ProductService.InsertProduct(product)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// GetManager 修改接口处理
func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id,err :=strconv.ParseInt(idString,10,16)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product,err:=p.ProductService.GetProductByID(id)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name:"product/manager.html",
		Data:iris.Map{
			"product":product,
		},
	}
}

func (p *ProductController) GetDelete() {
	idString:=p.Ctx.URLParam("id")
	id ,err := strconv.ParseInt(idString,10,64)
	if err !=nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOk:=p.ProductService.DeleteProductByID(id)
	if isOk{
		p.Ctx.Application().Logger().Debug("删除商品成功，ID为："+idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败，ID为："+idString)
	}
	p.Ctx.Redirect("/product/all")
}


