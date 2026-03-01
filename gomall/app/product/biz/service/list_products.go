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

package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ListProductsService struct {
	ctx context.Context
} 

func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

func convertProductToProto(v *model.Product, stock int64) *product.Product {
	p := &product.Product{
		Id:            uint32(v.ID),
		Name:          v.Name,
		Description:   v.Description,
		Picture:       v.Picture,
		Price:         v.Price,
		Sales:         uint32(v.Sales),
		Stock:         stock,
		DiscountType:  int32(v.DiscountType),
		DiscountValue: v.DiscountValue,
		DiscountLabel: v.GetDiscountLabel(),
		DiscountStatus: int32(v.GetDiscountStatus()),
		ActualPrice:   v.GetActualPrice(),
	}
	
	if v.DiscountStartTime != nil {
		p.DiscountStartTime = v.DiscountStartTime.Unix()
	}
	if v.DiscountEndTime != nil {
		p.DiscountEndTime = v.DiscountEndTime.Unix()
	}
	if v.OriginalPrice != nil {
		p.OriginalPrice = *v.OriginalPrice
	}
	
	return p
}

func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp = &product.ListProductsResp{}
	stockService := NewStockService(s.ctx)

	if req.DiscountFilter > 0 {
		products, err := model.GetProductsByDiscountFilter(mysql.DB, s.ctx, req.DiscountFilter)
		if err != nil {
			return nil, err
		}
		for _, v := range products {
			var stock int64 = 999
			stockData, _ := stockService.GetStock(uint32(v.ID))
			if stockData != nil {
				stock = stockData.Available
			}
			resp.Products = append(resp.Products, convertProductToProto(&v, stock))
		}
	} else if req.CategoryName != "" {
		c, err := model.GetProductsByCategoryName(mysql.DB, s.ctx, req.CategoryName)
		if err != nil {
			return nil, err
		}
		for _, v1 := range c {
			for _, v := range v1.Products {
				var stock int64 = 999
				stockData, _ := stockService.GetStock(uint32(v.ID))
				if stockData != nil {
					stock = stockData.Available
				}
				resp.Products = append(resp.Products, convertProductToProto(&v, stock))
			}
		}
	} else {
		products, err := model.GetAllProductsOrderByTime(mysql.DB, s.ctx)
		if err != nil {
			klog.Errorf("ListProducts: failed to get all products: %v", err)
			return nil, err
		}
		klog.Infof("ListProducts: found %d products from database", len(products))
		for _, v := range products {
			klog.Infof("ListProducts: Product ID=%d, Name=%s, Sales=%d, Price=%.2f, DiscountType=%d, DiscountStatus=%d", 
				v.ID, v.Name, v.Sales, v.Price, v.DiscountType, v.GetDiscountStatus())
			var stock int64 = 999
			stockData, _ := stockService.GetStock(uint32(v.ID))
			if stockData != nil {
				stock = stockData.Available
			}
			resp.Products = append(resp.Products, convertProductToProto(&v, stock))
		}
	}

	return resp, nil
}
