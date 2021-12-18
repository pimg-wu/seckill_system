package main

import (
	"fmt"
	"shopping_system/common"
	"shopping_system/rabbitmq"
	"shopping_system/repositories"
	"shopping_system/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)

	order := repositories.NewOrderManagerRepository("order", db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("product")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
