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
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type GetProductService struct {
	ctx context.Context
}

func NewGetProductService(ctx context.Context) *GetProductService {
	return &GetProductService{ctx: ctx}
}

func (s *GetProductService) Run(req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	if req.Id == 0 {
		return nil, kerrors.NewBizStatusError(40000, "product id is required")
	}

	p, err := model.NewProductQuery(s.ctx, mysql.DB).GetById(int(req.Id))
	if err != nil {
		return nil, err
	}
	
	var stock int64 = 999
	stockService := NewStockService(s.ctx)
	stockData, err := stockService.GetStock(uint32(p.ID))
	if err == nil && stockData != nil {
		stock = stockData.Available
	}
	
	protoProduct := &product.Product{
		Id:             uint32(p.ID),
		Picture:        p.Picture,
		Price:          p.Price,
		Description:    p.Description,
		Name:           p.Name,
		Sales:          uint32(p.Sales),
		Stock:          stock,
		DiscountType:   int32(p.DiscountType),
		DiscountValue:  p.DiscountValue,
		DiscountLabel:  p.GetDiscountLabel(),
		DiscountStatus: int32(p.GetDiscountStatus()),
		ActualPrice:    p.GetActualPrice(),
	}
	
	if p.DiscountStartTime != nil {
		protoProduct.DiscountStartTime = p.DiscountStartTime.Unix()
	}
	if p.DiscountEndTime != nil {
		protoProduct.DiscountEndTime = p.DiscountEndTime.Unix()
	}
	if p.OriginalPrice != nil {
		protoProduct.OriginalPrice = *p.OriginalPrice
	}
	
	return &product.GetProductResp{
		Product: protoProduct,
	}, err
}
