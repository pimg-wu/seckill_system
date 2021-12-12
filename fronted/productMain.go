package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	app.HandleDir("/public", "./web/public")
	app.HandleDir("/html", "./web/htmlProductShow")

	app.Run(iris.Addr("0.0.0.0:80"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
