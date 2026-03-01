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
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"gorm.io/gorm"
)

type AdminProductsService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAdminProductsService(Context context.Context, RequestContext *app.RequestContext) *AdminProductsService {
	return &AdminProductsService{RequestContext: RequestContext, Context: Context}
}

type Product struct {
	Id          uint32         `gorm:"column:id;primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Picture     string         `json:"picture"`
	Price       float32        `json:"price"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "product"
}

func (s *AdminProductsService) Run() (res map[string]any, err error) {
	db, err := utils.GetProductDB()
	if err != nil {
		return hertzutils.H{
			"title":    "商品管理",
			"products": []Product{},
		}, nil
	}

	var products []Product
	err = db.Find(&products).Error
	if err != nil {
		return hertzutils.H{
			"title":    "商品管理",
			"products": []Product{},
		}, nil
	}

	return hertzutils.H{
		"title":    "商品管理",
		"products": products,
	}, nil
}
