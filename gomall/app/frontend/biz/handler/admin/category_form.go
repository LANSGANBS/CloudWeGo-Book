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

func CategoryForm(ctx context.Context, c *app.RequestContext) {
	categoryID := c.Query("id")

	resp := hertzutils.H{
		"title":      "分类表单",
		"categoryID": categoryID,
		"category":   nil,
	}

	if categoryID != "" {
		id, err := strconv.Atoi(categoryID)
		if err != nil {
			hlog.CtxErrorf(ctx, "invalid category id: %v", err)
		} else {
			db, err := utils.GetProductDB()
			if err != nil {
				hlog.CtxErrorf(ctx, "failed to get db: %v", err)
			} else {
				var cat Category
				err = db.First(&cat, id).Error
				if err != nil {
					hlog.CtxWarnf(ctx, "category not found: %v", err)
				} else {
					resp["category"] = cat
				}
			}
		}
	}

	c.HTML(consts.StatusOK, "admin-category-form", utils.WarpResponse(ctx, c, resp))
}
