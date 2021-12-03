package main

import (
	"context"
	"log"
	"shopping_system/backend/web/controllers"
	"shopping_system/common"
	"shopping_system/repositories"
	"shopping_system/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	//创建iris实例
	app := iris.New()

	app.Logger().SetLevel("debug") //mvc模式下提示错误

	//注册模板
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	//设置模板目标
	// app.StaticWeb("/assets", "./backend/web/assets")
	app.HandleDir("/assets", "./web/assets")

	//出现异常跳转指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	//启动服务
	app.Run(iris.Addr("localhost:8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
