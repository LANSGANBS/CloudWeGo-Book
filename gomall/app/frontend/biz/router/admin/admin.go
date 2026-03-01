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

package admin

import (
	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/handler/admin"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// Register register admin routes
func Register(r *server.Hertz) {
	adminGroup := r.Group("/admin", rootMw()...)

	// 管理首页
	adminGroup.GET("/", append(_adminMw(), admin.Dashboard)...)

	// 商品管理
	adminGroup.GET("/products", append(_adminMw(), admin.Products)...)
	adminGroup.GET("/products/form", append(_adminMw(), admin.ProductForm)...)
	adminGroup.POST("/products/save", append(_adminMw(), admin.SaveProduct)...)
	adminGroup.DELETE("/products/delete", append(_adminMw(), admin.DeleteProduct)...)

	// 分类管理
	adminGroup.GET("/categories", append(_adminMw(), admin.Categories)...)
	adminGroup.GET("/categories/form", append(_adminMw(), admin.CategoryForm)...)
	adminGroup.POST("/categories/save", append(_adminMw(), admin.SaveCategory)...)
	adminGroup.DELETE("/categories/delete", append(_adminMw(), admin.DeleteCategory)...)

	// 订单管理
	adminGroup.GET("/orders", append(_adminMw(), admin.Orders)...)

	// 图片管理
	adminGroup.GET("/images", append(_adminMw(), admin.ListImages)...)
}
