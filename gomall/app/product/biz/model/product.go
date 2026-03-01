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

package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	DiscountTypeNone     int8 = 0
	DiscountTypeRate     int8 = 1
	DiscountTypeFixed    int8 = 2
)

type Product struct {
	ID                int            `gorm:"primarykey;column:id"`
	CreatedAt         time.Time      `gorm:"column:created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at"`
	Name              string         `json:"name" gorm:"column:name"`
	Description       string         `json:"description" gorm:"column:description"`
	Picture           string         `json:"picture" gorm:"column:picture"`
	Price             float32        `json:"price" gorm:"column:price"`
	DiscountType      int8           `json:"discount_type" gorm:"column:discount_type;default:0"`
	DiscountValue     float32        `json:"discount_value" gorm:"column:discount_value;default:0"`
	DiscountStartTime *time.Time     `json:"discount_start_time" gorm:"column:discount_start_time"`
	DiscountEndTime   *time.Time     `json:"discount_end_time" gorm:"column:discount_end_time"`
	OriginalPrice     *float32       `json:"original_price" gorm:"column:original_price"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Sales             int64          `json:"sales" gorm:"column:sales;default:0"`
	Categories        []Category     `json:"categories" gorm:"many2many:product_category"`
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

func (p *Product) IsFlashSale() bool {
	return p.DiscountType != DiscountTypeNone && p.DiscountStartTime != nil && p.DiscountEndTime != nil
}

func (p *Product) IsNormalDiscount() bool {
	return p.DiscountType != DiscountTypeNone && p.DiscountStartTime == nil
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

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (p ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Where(&Product{ID: productId}).First(&product).Error
	return
}

func NewProductQuery(ctx context.Context, db *gorm.DB) ProductQuery {
	return ProductQuery{ctx: ctx, db: db}
}

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
}

func (c CachedProductQuery) GetById(productId int) (product Product, err error) {
	cacheKey := fmt.Sprintf("%s_%s_%d", c.prefix, "product_by_id", productId)
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cacheKey)

	err = func() error {
		err1 := cachedResult.Err()
		if err1 != nil {
			return err1
		}
		cachedResultByte, err2 := cachedResult.Bytes()
		if err2 != nil {
			return err2
		}
		err3 := json.Unmarshal(cachedResultByte, &product)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return Product{}, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return
}

func NewCachedProductQuery(pq ProductQuery, cacheClient *redis.Client) CachedProductQuery {
	return CachedProductQuery{productQuery: pq, cacheClient: cacheClient, prefix: "cloudwego_shop"}
}

func GetProductById(db *gorm.DB, ctx context.Context, productId int) (product Product, err error) {
	err = db.WithContext(ctx).Model(&Product{}).Select("id, created_at, updated_at, name, description, picture, price, discount_type, discount_value, discount_start_time, discount_end_time, original_price, deleted_at, sales").Where(&Product{ID: productId}).First(&product).Error
	return product, err
}

func SearchProduct(db *gorm.DB, ctx context.Context, q string) (product []*Product, err error) {
	err = db.WithContext(ctx).Model(&Product{}).Select("id, created_at, updated_at, name, description, picture, price, discount_type, discount_value, discount_start_time, discount_end_time, original_price, deleted_at, sales").Find(&product, "name like ? or description like ?", "%"+q+"%", "%"+q+"%").Error
	return product, err
}

func IncrementSales(db *gorm.DB, ctx context.Context, productId int, quantity int) error {
	return db.WithContext(ctx).Model(&Product{}).Where("id = ?", productId).
		UpdateColumn("sales", gorm.Expr("sales + ?", quantity)).Error
}
