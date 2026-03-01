// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type ProductCatalogServiceImpl struct{}

func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp, err = service.NewListProductsService(ctx).Run(req)

	return resp, err
}

func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	resp, err = service.NewGetProductService(ctx).Run(req)

	return resp, err
}

func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	resp, err = service.NewSearchProductsService(ctx).Run(req)

	return resp, err
}

func (s *ProductCatalogServiceImpl) IncrementSales(ctx context.Context, req *product.IncrementSalesReq) (resp *product.IncrementSalesResp, err error) {
	resp, err = service.NewIncrementSalesService(ctx).Run(req)

	return resp, err
}

func (s *ProductCatalogServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	resp = &product.CreateProductResp{}
	return resp, err
}

func (s *ProductCatalogServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	resp = &product.UpdateProductResp{}
	return resp, err
}

func (s *ProductCatalogServiceImpl) DeductStock(ctx context.Context, req *product.DeductStockReq) (resp *product.DeductStockResp, err error) {
	resp, err = service.NewDeductStockService(ctx).Run(req)
	return resp, err
}

func (s *ProductCatalogServiceImpl) RestoreStock(ctx context.Context, req *product.RestoreStockReq) (resp *product.RestoreStockResp, err error) {
	resp, err = service.NewRestoreStockService(ctx).Run(req)
	return resp, err
}

func (s *ProductCatalogServiceImpl) GetStock(ctx context.Context, req *product.GetStockReq) (resp *product.GetStockResp, err error) {
	resp, err = service.NewGetStockService(ctx).Run(req)
	return resp, err
}

func (s *ProductCatalogServiceImpl) SetDiscount(ctx context.Context, req *product.SetDiscountReq) (resp *product.SetDiscountResp, err error) {
	resp, err = service.NewSetDiscountService(ctx).Run(req)
	return resp, err
}

func (s *ProductCatalogServiceImpl) CancelDiscount(ctx context.Context, req *product.CancelDiscountReq) (resp *product.CancelDiscountResp, err error) {
	resp, err = service.NewCancelDiscountService(ctx).Run(req)
	return resp, err
}

func (s *ProductCatalogServiceImpl) GetProductPriceHistory(ctx context.Context, req *product.GetProductPriceHistoryReq) (resp *product.GetProductPriceHistoryResp, err error) {
	resp, err = service.NewGetProductPriceHistoryService(ctx).Run(req)
	return resp, err
}
