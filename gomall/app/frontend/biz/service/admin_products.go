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
	"fmt"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/hertz/pkg/app"
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"gorm.io/gorm"
)

type AdminProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAdminProductsService(Context context.Context, RequestContext *app.RequestContext) *AdminProductsService {
	return &AdminProductsService{RequestContext: RequestContext, Context: Context}
}

const (
	DiscountTypeNone  int8 = 0
	DiscountTypeRate  int8 = 1
	DiscountTypeFixed int8 = 2
)

type Product struct {
	Id                uint32         `gorm:"column:id;primaryKey" json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Picture           string         `json:"picture"`
	Price             float32        `json:"price"`
	DiscountType      int8           `gorm:"column:discount_type;default:0" json:"discount_type"`
	DiscountValue     float32        `gorm:"column:discount_value;default:0" json:"discount_value"`
	DiscountStartTime *time.Time     `gorm:"column:discount_start_time" json:"discount_start_time"`
	DiscountEndTime   *time.Time     `gorm:"column:discount_end_time" json:"discount_end_time"`
	OriginalPrice     *float32       `gorm:"column:original_price" json:"original_price"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "product"
}

func (p *Product) IsOnDiscount() bool {
	if p.DiscountType == DiscountTypeNone {
		return false
	}
	if p.DiscountStartTime == nil && p.DiscountEndTime == nil {
		return true
	}
	now := time.Now()
	if p.DiscountStartTime != nil && now.Before(*p.DiscountStartTime) {
		return false
	}
	if p.DiscountEndTime != nil && now.After(*p.DiscountEndTime) {
		return false
	}
	return true
}

func (p *Product) GetDiscountStatus() int {
	if p.DiscountType == DiscountTypeNone {
		return 0
	}
	if p.DiscountStartTime != nil && p.DiscountEndTime != nil {
		now := time.Now()
		if now.Before(*p.DiscountStartTime) {
			return 1
		}
		if now.After(*p.DiscountEndTime) {
			return 3
		}
		return 2
	}
	return 4
}

func (p *Product) GetActualPrice() float32 {
	if !p.IsOnDiscount() {
		return p.Price
	}
	switch p.DiscountType {
	case DiscountTypeRate:
		return p.Price * p.DiscountValue
	case DiscountTypeFixed:
		actualPrice := p.Price - p.DiscountValue
		if actualPrice < 0 {
			return 0
		}
		return actualPrice
	default:
		return p.Price
	}
}

func (p *Product) GetDiscountLabel() string {
	if !p.IsOnDiscount() {
		return ""
	}
	switch p.DiscountType {
	case DiscountTypeRate:
		discountRate := int(p.DiscountValue * 10)
		return fmt.Sprintf("%d折", discountRate)
	case DiscountTypeFixed:
		return fmt.Sprintf("减%.0f元", p.DiscountValue)
	default:
		return ""
	}
}

func convertRPCProductToProduct(p *product.Product) *Product {
	prod := &Product{
		Id:            p.Id,
		Name:          p.Name,
		Description:   p.Description,
		Picture:       p.Picture,
		Price:         p.Price,
		DiscountType:  int8(p.DiscountType),
		DiscountValue: p.DiscountValue,
	}
	
	if p.DiscountStartTime > 0 {
		t := time.Unix(p.DiscountStartTime, 0)
		prod.DiscountStartTime = &t
	}
	if p.DiscountEndTime > 0 {
		t := time.Unix(p.DiscountEndTime, 0)
		prod.DiscountEndTime = &t
	}
	if p.OriginalPrice > 0 {
		prod.OriginalPrice = &p.OriginalPrice
	}
	
	return prod
}

func (s *AdminProductsService) Run() (res map[string]any, err error) {
	ctx := s.Context
	
	resp, err := rpc.ProductClient.ListProducts(ctx, &product.ListProductsReq{})
	if err != nil {
		hlog.CtxErrorf(ctx, "AdminProducts: failed to call ListProducts RPC: %v", err)
		return hertzutils.H{
			"title":    "商品管理",
			"products": []Product{},
		}, nil
	}

	hlog.CtxInfof(ctx, "AdminProducts: found %d products from RPC", len(resp.Products))

	products := make([]*Product, 0, len(resp.Products))
	for _, p := range resp.Products {
		products = append(products, convertRPCProductToProduct(p))
	}

	hlog.CtxInfof(ctx, "AdminProducts: converted %d products", len(products))

	return hertzutils.H{
		"title":    "商品管理",
		"products": products,
	}, nil
}
