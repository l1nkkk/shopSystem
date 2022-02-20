package controllers

import (
	"github.com/kataras/iris/v12/v12/mvc"
	"github.com/l1nkkk/shopSystem/demo/irisdemo/repositories"
	"github.com/l1nkkk/shopSystem/demo/irisdemo/services"
)
type MovieController struct {

}

func (c *MovieController) Get() mvc.View{
	// 1. 创建数据库操作对象
	movieRepository := repositories.NewMovieManager()
	// 2. 建立 service,处理业务
	movieService := services.NewMovieServiceManger(movieRepository)
	res := movieService.ShowMovieName()

	// 3. 渲染
	return mvc.View{
		Name: "movie/index.html",
		Data: res,
	}
}