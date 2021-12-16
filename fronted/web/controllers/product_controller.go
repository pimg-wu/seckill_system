package controllers

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	datamodels "shopping_system/dataModels"
	"shopping_system/rabbitmq"
	"shopping_system/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMq       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

var (
	htmlOutPath  = "./web/htmlProductShow/" //生成的html目录
	templatePath = "./web/views/template/"  //静态文件模板
)

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//获取模板
	contenstTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	product, err := p.ProductService.GetProductById(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	generateStaticHtml(p.Ctx, contenstTmp, fileName, product)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	//判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName)
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
	}
	//生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	defer file.Close()
	template.Execute(file, &product) //模板渲染
}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductById(1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// func (p *ProductController) GetOrder() mvc.View {
// 	productString := p.Ctx.URLParam("productID")
// 	userString := p.Ctx.GetCookie("uid")
// 	productID, err := strconv.Atoi(productString)
// 	if err != nil {
// 		p.Ctx.Application().Logger().Debug(err)
// 	}
// 	product, err := p.ProductService.GetProductById(int64(productID))
// 	if err != nil {
// 		p.Ctx.Application().Logger().Debug(err)
// 	}
// 	var orderID int64
// 	showMessage := "抢购失败！"
//
// 	//判断商品数量是否满足需求
// 	if product.ProductNum > 0 {
// 		//扣除商品数量
// 		product.ProductNum -= 1
// 		err := p.ProductService.UpdateProduct(product) //高并发下超卖！
// 		if err != nil {
// 			p.Ctx.Application().Logger().Debug(err)
// 		}
// 		//创建订单
// 		userID, err := strconv.Atoi(userString)
// 		if err != nil {
// 			p.Ctx.Application().Logger().Debug(err)
// 		}
// 		order := &datamodels.Order{
// 			UserId:      int64(userID),
// 			ProductId:   int64(productID),
// 			OrderStatus: datamodels.OrderSuccess,
// 		}
// 		//新建订单
// 		orderID, err = p.OrderService.InsertOrder(order)
// 		if err != nil {
// 			p.Ctx.Application().Logger().Debug(err)
// 		} else {
// 			showMessage = "抢购成功！"
// 		}
// 	}
//
// 	return mvc.View{
// 		Layout: "shared/productLayout.html",
// 		Name:   "product/result.html",
// 		Data: iris.Map{
// 			"orderID":     orderID,
// 			"showMessage": showMessage,
// 		},
// 	}
// }

func (p *ProductController) GetOrder() []byte {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//创建消息体
	message := datamodels.NewMessage(userID, productID)
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err = p.RabbitMq.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")
}
