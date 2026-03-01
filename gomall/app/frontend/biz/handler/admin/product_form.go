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
	"context"
	"strconv"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func ProductForm(ctx context.Context, c *app.RequestContext) {
	productID := c.Query("id")

	db, err := utils.GetProductDB()
	if err != nil {
		c.HTML(consts.StatusInternalServerError, "error", hertzutils.H{
			"error": "数据库连接失败: " + err.Error(),
		})
		return
	}

	var categories []Category
	err = db.Find(&categories).Error
	if err != nil {
		c.HTML(consts.StatusInternalServerError, "error", hertzutils.H{
			"error": "获取分类列表失败: " + err.Error(),
		})
		return
	}

	resp := hertzutils.H{
		"title":      "商品表单",
		"productID":  productID,
		"categories": categories,
		"product":    nil,
		"productCategory": "",
	}

	// 编辑模式：加载现有商品数据
	if productID != "" {
		id, err := strconv.Atoi(productID)
		if err != nil {
			hlog.CtxErrorf(ctx, "invalid product id: %v", err)
		} else {
			var product Product
			err = db.Preload("Categories").First(&product, id).Error
			if err != nil {
				hlog.CtxWarnf(ctx, "product not found: %v", err)
			} else {
				// 加载库存数据
				var stock Stock
				stockErr := db.Where("product_id = ?", id).First(&stock).Error
				if stockErr == nil {
					product.Stock = stock.Available
				} else {
					product.Stock = 100 // 默认库存
				}
				resp["product"] = product
				// 设置当前分类
				if len(product.Categories) > 0 {
					resp["productCategory"] = product.Categories[0].Name
				}
			}
		}
	}

	c.HTML(consts.StatusOK, "admin-product-form", utils.WarpResponse(ctx, c, resp))
}
