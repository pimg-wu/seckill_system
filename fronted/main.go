package main

import (
	"context"
	"shopping_system/common"
	"shopping_system/fronted/middleware"
	"shopping_system/fronted/web/controllers"
	"shopping_system/repositories"
	"shopping_system/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	tmplate := iris.HTML("./web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	app.HandleDir("/public", "./web/public")
	app.HandleDir("/html", "./web/htmlProductShow")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {

	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	app.Run(iris.Addr("0.0.0.0:8082"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
