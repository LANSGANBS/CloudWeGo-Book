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
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

type Product struct {
	ID          int            `gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Picture     string         `json:"picture"`
	Price       float32        `json:"price"`
	Stock       int64          `json:"stock" gorm:"-"`
	Categories  []Category     `json:"categories" gorm:"many2many:product_category"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "product"
}

type Category struct {
	ID          int            `gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Category) TableName() string {
	return "category"
}

func SaveProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	picture := c.PostForm("picture")
	category := c.PostForm("category")
	stockStr := c.PostForm("stock")

	if name == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "商品名称不能为空",
		})
		return
	}
	if description == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "商品描述不能为空",
		})
		return
	}
	if picture == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "商品图片不能为空",
		})
		return
	}
	if category == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "请选择商品分类",
		})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "价格格式错误",
		})
		return
	}
	if price <= 0 {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "价格必须大于0",
		})
		return
	}

	var stock int64 = 100
	if stockStr != "" {
		stock, err = strconv.ParseInt(stockStr, 10, 64)
		if err != nil {
			c.JSON(consts.StatusBadRequest, hertzutils.H{
				"success": false,
				"message": "库存格式错误",
			})
			return
		}
		if stock < 0 {
			c.JSON(consts.StatusBadRequest, hertzutils.H{
				"success": false,
				"message": "库存不能为负数",
			})
			return
		}
	}

	db, err := utils.GetProductDB()
	if err != nil {
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "数据库连接失败: " + err.Error(),
		})
		return
	}

	var cat Category
	err = db.Where("name = ?", category).First(&cat).Error
	if err != nil {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类不存在: " + category,
		})
		return
	}

	// 编辑模式
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(consts.StatusBadRequest, hertzutils.H{
				"success": false,
				"message": "商品ID格式错误",
			})
			return
		}

		var existingProduct Product
		err = db.First(&existingProduct, id).Error
		if err != nil {
			c.JSON(consts.StatusNotFound, hertzutils.H{
				"success": false,
				"message": "商品不存在",
			})
			return
		}

		// 更新商品基本信息
		existingProduct.Name = name
		existingProduct.Description = description
		existingProduct.Price = float32(price)
		existingProduct.Picture = picture
		existingProduct.UpdatedAt = time.Now()

		// 更新商品分类关联
		err = db.Model(&existingProduct).Association("Categories").Replace(&cat)
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update product categories: %v", err)
			c.JSON(consts.StatusInternalServerError, hertzutils.H{
				"success": false,
				"message": "更新分类失败: " + err.Error(),
			})
			return
		}

		err = db.Save(&existingProduct).Error
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update product: %v", err)
			c.JSON(consts.StatusInternalServerError, hertzutils.H{
				"success": false,
				"message": "更新失败: " + err.Error(),
			})
			return
		}

		// 更新库存
		err = updateOrCreateStock(db, uint32(id), stock)
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update stock: %v", err)
		}

		c.JSON(consts.StatusOK, hertzutils.H{
			"success": true,
			"message": "商品更新成功",
		})
		return
	}

	// 新增模式
	now := time.Now()
	product := Product{
		Name:        name,
		Description: description,
		Price:       float32(price),
		Picture:     picture,
		Categories:  []Category{cat},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = db.Create(&product).Error
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to create product: %v", err)
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "保存失败: " + err.Error(),
		})
		return
	}

	// 创建库存记录
	err = updateOrCreateStock(db, uint32(product.ID), stock)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to create stock: %v", err)
	}

	c.JSON(consts.StatusOK, hertzutils.H{
		"success": true,
		"message": "商品创建成功",
	})
}

func updateOrCreateStock(db *gorm.DB, productId uint32, quantity int64) error {
	var existingStock Stock
	err := db.Where("product_id = ?", productId).First(&existingStock).Error
	if err == gorm.ErrRecordNotFound {
		stock := Stock{
			ProductId:   productId,
			Quantity:    quantity,
			Available:   quantity,
			Reserved:    0,
			MinStock:    10,
			MaxStock:    1000,
			SafetyStock: 20,
			Status:      1,
		}
		return db.Create(&stock).Error
	} else if err != nil {
		return err
	}
	
	return db.Model(&existingStock).Updates(map[string]interface{}{
		"quantity":  quantity,
		"available": quantity - existingStock.Reserved,
	}).Error
}

type Stock struct {
	ID          uint32         `gorm:"primaryKey;autoIncrement"`
	ProductId   uint32         `gorm:"uniqueIndex;not null"`
	Quantity    int64          `gorm:"default:0"`
	Reserved    int64          `gorm:"default:0"`
	Available   int64          `gorm:"default:0"`
	MinStock    int64          `gorm:"default:10"`
	MaxStock    int64          `gorm:"default:1000"`
	SafetyStock int64          `gorm:"default:20"`
	Status      int8           `gorm:"default:1"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Stock) TableName() string {
	return "stock"
}
