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

type CategoryForSave struct {
	ID          int            `gorm:"primaryKey"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (CategoryForSave) TableName() string {
	return "category"
}

func SaveCategory(ctx context.Context, c *app.RequestContext) {
	idStr := c.PostForm("id")
	name := c.PostForm("name")
	description := c.PostForm("description")

	if name == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类名称不能为空",
		})
		return
	}
	if description == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类描述不能为空",
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

	// 编辑模式
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(consts.StatusBadRequest, hertzutils.H{
				"success": false,
				"message": "分类ID格式错误",
			})
			return
		}

		var existingCat CategoryForSave
		err = db.First(&existingCat, id).Error
		if err != nil {
			c.JSON(consts.StatusNotFound, hertzutils.H{
				"success": false,
				"message": "分类不存在",
			})
			return
		}

		// 检查名称是否与其他分类重复
		var duplicateCat CategoryForSave
		err = db.Where("name = ? AND id != ?", name, id).First(&duplicateCat).Error
		if err == nil {
			c.JSON(consts.StatusBadRequest, hertzutils.H{
				"success": false,
				"message": "分类名称已被其他分类使用",
			})
			return
		}

		existingCat.Name = name
		existingCat.Description = description
		existingCat.UpdatedAt = time.Now()

		err = db.Save(&existingCat).Error
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update category: %v", err)
			c.JSON(consts.StatusInternalServerError, hertzutils.H{
				"success": false,
				"message": "更新失败: " + err.Error(),
			})
			return
		}

		c.JSON(consts.StatusOK, hertzutils.H{
			"success": true,
			"message": "分类更新成功",
		})
		return
	}

	// 新增模式
	var existingCat CategoryForSave
	err = db.Where("name = ?", name).First(&existingCat).Error
	if err == nil {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "分类名称已存在",
		})
		return
	}

	now := time.Now()
	category := CategoryForSave{
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = db.Create(&category).Error
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to create category: %v", err)
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "保存失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, hertzutils.H{
		"success": true,
		"message": "分类创建成功",
	})
}
