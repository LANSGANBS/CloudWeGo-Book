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

	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"product" gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

func GetProductsByCategoryName(db *gorm.DB, ctx context.Context, name string) (category []Category, err error) {
	err = db.WithContext(ctx).Model(&Category{}).Where(&Category{Name: name}).Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, created_at, updated_at, name, description, picture, price, discount_type, discount_value, discount_start_time, discount_end_time, original_price, deleted_at, sales").Order("updated_at DESC")
	}).Find(&category).Error
	return category, err
}

func GetAllProductsOrderByTime(db *gorm.DB, ctx context.Context) (products []Product, err error) {
	err = db.WithContext(ctx).Model(&Product{}).Select("id, created_at, updated_at, name, description, picture, price, discount_type, discount_value, discount_start_time, discount_end_time, original_price, deleted_at, sales").Order("updated_at DESC").Find(&products).Error
	return products, err
}

func GetProductsByDiscountFilter(db *gorm.DB, ctx context.Context, discountFilter int32) (products []Product, err error) {
	query := db.WithContext(ctx).Model(&Product{}).Select("id, created_at, updated_at, name, description, picture, price, discount_type, discount_value, discount_start_time, discount_end_time, original_price, deleted_at, sales")
	
	switch discountFilter {
	case 1:
		query = query.Where("discount_type != 0")
	case 2:
		query = query.Where("discount_type != 0 AND discount_start_time IS NULL")
	case 3:
		query = query.Where("discount_type != 0 AND discount_start_time IS NOT NULL AND discount_end_time IS NOT NULL")
	}
	
	err = query.Order("updated_at DESC").Find(&products).Error
	return products, err
}
