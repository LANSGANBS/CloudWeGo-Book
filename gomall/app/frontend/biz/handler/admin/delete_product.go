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
	hertzutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func DeleteProduct(ctx context.Context, c *app.RequestContext) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "商品ID不能为空",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(consts.StatusBadRequest, hertzutils.H{
			"success": false,
			"message": "商品ID格式错误",
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

	err = db.Delete(&Product{}, id).Error
	if err != nil {
		c.JSON(consts.StatusInternalServerError, hertzutils.H{
			"success": false,
			"message": "删除失败: " + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, hertzutils.H{
		"success": true,
		"message": "商品删除成功",
	})
}
