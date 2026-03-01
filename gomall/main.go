package main

import (
	product "github.com/cloudwego/biz-demo/gomall/kitex_gen/product/productcatalogservice"
	"log"
)

func main() {
	svr := product.NewServer(new(ProductCatalogServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
