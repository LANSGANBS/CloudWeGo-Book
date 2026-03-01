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

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"gorm.io/gorm"
)

type AdminCategoriesService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAdminCategoriesService(Context context.Context, RequestContext *app.RequestContext) *AdminCategoriesService {
	return &AdminCategoriesService{RequestContext: RequestContext, Context: Context}
}

type Category struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ProductNum  int            `gorm:"-" json:"product_num"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "category"
}

func (s *AdminCategoriesService) Run() (res map[string]any, err error) {
	db, err := utils.GetProductDB()
	if err != nil {
		hlog.CtxErrorf(s.Context, "failed to get product db: %v", err)
		return hertzutils.H{
			"title":                  "分类管理",
			"categories":             []Category{},
			"categoriesWithProducts": 0,
			"totalProducts":          0,
		}, nil
	}

	var categories []Category
	err = db.Find(&categories).Error
	if err != nil {
		hlog.CtxErrorf(s.Context, "failed to query categories: %v", err)
		return hertzutils.H{
			"title":                  "分类管理",
			"categories":             []Category{},
			"categoriesWithProducts": 0,
			"totalProducts":          0,
		}, nil
	}

	categoriesWithProducts := 0
	totalProducts := 0

	for i := range categories {
		var count int64
		err = db.Table("product_category").
			Joins("JOIN product ON product_category.product_id = product.id").
			Where("product_category.category_id = ? AND product.deleted_at IS NULL", categories[i].ID).
			Count(&count).Error
		if err != nil {
			hlog.CtxWarnf(s.Context, "failed to count products for category %d: %v", categories[i].ID, err)
			count = 0
		}
		categories[i].ProductNum = int(count)
		totalProducts += int(count)
		if count > 0 {
			categoriesWithProducts++
		}
	}

	return hertzutils.H{
		"title":                  "分类管理",
		"categories":             categories,
		"categoriesWithProducts": categoriesWithProducts,
		"totalProducts":          totalProducts,
	}, nil
}
