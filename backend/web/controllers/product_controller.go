package controllers

import (
	"shopping_system/common"
	datamodels "shopping_system/dataModels"
	"shopping_system/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}

func (p *ProductController) GetAll() mvc.View {
	productArray, _ := p.ProductService.GetAllProduct()
	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": productArray,
		},
	}
}

//修改商品
func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductById(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOK := p.ProductService.DeleteProductById(id)
	if isOK {
		p.Ctx.Application().Logger().Debug("删除商品成功， ID为：" + idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败， ID为：" + idString)
	}
	p.Ctx.Redirect("/product/all")
}
