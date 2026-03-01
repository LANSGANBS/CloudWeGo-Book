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
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service/rag"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service/vector"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type SearchProductsService struct {
	ctx context.Context
}

func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	ragSvc := rag.GetRAGService()
	stockService := NewStockService(s.ctx)

	results, err := ragSvc.Search(s.ctx, req.Query, 10)
	if err != nil {
		klog.Warnf("RAG search failed, falling back to SQL search: %v", err)
		return s.fallbackSearch(req)
	}

	if len(results) == 0 {
		klog.Info("RAG search returned no results, falling back to SQL search")
		return s.fallbackSearch(req)
	}

	var productResults []*product.Product
	for _, r := range results {
		p, err := model.GetProductById(mysql.DB, s.ctx, int(r.ProductID))
		if err != nil {
			klog.Warnf("Failed to get product %d: %v", r.ProductID, err)
			continue
		}
		var stock int64 = 999
		stockData, _ := stockService.GetStock(uint32(p.ID))
		if stockData != nil {
			stock = stockData.Available
		}
		protoProduct := &product.Product{
			Id:             uint32(p.ID),
			Name:           p.Name,
			Description:    p.Description,
			Picture:        p.Picture,
			Price:          p.Price,
			Stock:          stock,
			Sales:          uint32(p.Sales),
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
		productResults = append(productResults, protoProduct)
	}

	klog.Infof("RAG search returned %d products for query: %s", len(productResults), req.Query)
	return &product.SearchProductsResp{Results: productResults}, nil
}

func (s *SearchProductsService) fallbackSearch(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	p, err := model.SearchProduct(mysql.DB, s.ctx, req.Query)
	if err != nil {
		return nil, err
	}

	stockService := NewStockService(s.ctx)
	var results []*product.Product
	for _, v := range p {
		var stock int64 = 999
		stockData, _ := stockService.GetStock(uint32(v.ID))
		if stockData != nil {
			stock = stockData.Available
		}
		protoProduct := &product.Product{
			Id:             uint32(v.ID),
			Name:           v.Name,
			Description:    v.Description,
			Picture:        v.Picture,
			Price:          v.Price,
			Stock:          stock,
			Sales:          uint32(v.Sales),
			DiscountType:   int32(v.DiscountType),
			DiscountValue:  v.DiscountValue,
			DiscountLabel:  v.GetDiscountLabel(),
			DiscountStatus: int32(v.GetDiscountStatus()),
			ActualPrice:    v.GetActualPrice(),
		}
		if v.DiscountStartTime != nil {
			protoProduct.DiscountStartTime = v.DiscountStartTime.Unix()
		}
		if v.DiscountEndTime != nil {
			protoProduct.DiscountEndTime = v.DiscountEndTime.Unix()
		}
		if v.OriginalPrice != nil {
			protoProduct.OriginalPrice = *v.OriginalPrice
		}
		results = append(results, protoProduct)
	}
	return &product.SearchProductsResp{Results: results}, nil
}

func (s *SearchProductsService) SearchWithScores(req *product.SearchProductsReq) ([]*vector.SearchResult, error) {
	ragSvc := rag.GetRAGService()
	return ragSvc.Search(s.ctx, req.Query, 10)
}
