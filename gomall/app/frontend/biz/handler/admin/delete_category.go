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
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

func DeleteCategory(ctx context.Context, c *app.RequestContext) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类ID不能为空",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类ID格式错误",
		})
		return
	}

	db, err := utils.GetProductDB()
	if err != nil {
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "数据库连接失败: " + err.Error(),
		})
		return
	}

	var count int64
	db.Table("product_category").
		Joins("JOIN product ON product_category.product_id = product.id").
		Where("product_category.category_id = ? AND product.deleted_at IS NULL", id).
		Count(&count)
	if count > 0 {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "该分类下还有商品，无法删除",
		})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Unscoped().Where("deleted_at IS NOT NULL").Delete(&Category{})
		if err := tx.Delete(&Category{}, id).Error; err != nil {
			return err
		}
		var validCategories []Category
		if err := tx.Where("deleted_at IS NULL").Order("id").Find(&validCategories).Error; err != nil {
			return err
		}
		idMap := make(map[int]int)
		for idx, cat := range validCategories {
			idMap[cat.ID] = idx + 1
		}
		if err := tx.Model(&Category{}).Where("id = ?", id).Update("deleted_at", gorm.DeletedAt{Time: time.Now()}).Error; err != nil {
			return err
		}
		for idx, cat := range validCategories {
			newID := idx + 1
			if err := tx.Model(&Category{}).Where("id = ?", cat.ID).Update("id", newID).Error; err != nil {
				return err
			}
		}
		for _, cat := range validCategories {
			if err := tx.Exec("UPDATE product_category SET category_id = ? WHERE category_id = ?", idMap[cat.ID]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, hertzutils.H{
		"success": true,
		"message": "分类删除成功",
	})
}
