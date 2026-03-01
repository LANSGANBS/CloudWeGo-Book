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

package utils

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	frontendutils "github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SendErrResponse  pack error response
func SendErrResponse(ctx context.Context, c *app.RequestContext, code int, err error) {
	// todo edit custom code
	c.String(code, err.Error())
}

// SendSuccessResponse  pack success response
func SendSuccessResponse(ctx context.Context, c *app.RequestContext, code int, data interface{}) {
	// todo edit custom code
	c.JSON(code, data)
}

// Category 分类结构
type Category struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "category"
}

var (
	productDB     *gorm.DB
	productDBOnce sync.Once
	productDBErr  error
)

// GetProductDB 获取产品数据库连接（单例模式，导出函数）
func GetProductDB() (*gorm.DB, error) {
	productDBOnce.Do(func() {
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mysqlHost := os.Getenv("MYSQL_HOST")

		if mysqlUser == "" {
			mysqlUser = "root"
		}
		if mysqlPassword == "" {
			mysqlPassword = "123456"
		}
		if mysqlHost == "" {
			mysqlHost = "127.0.0.1"
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/product?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlUser, mysqlPassword, mysqlHost)

		productDB, productDBErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if productDBErr != nil {
			hlog.Errorf("failed to connect to product database: %v", productDBErr)
		}
	})
	return productDB, productDBErr
}

func WarpResponse(ctx context.Context, c *app.RequestContext, content map[string]any) map[string]any {
	var cartNum int
	userId := frontendutils.GetUserIdFromCtx(ctx)
	cartResp, _ := rpc.CartClient.GetCart(ctx, &cart.GetCartReq{UserId: userId})
	if cartResp != nil && cartResp.Cart != nil {
		cartNum = len(cartResp.Cart.Items)
	}
	content["user_id"] = ctx.Value(frontendutils.UserIdKey)
	content["cart_num"] = cartNum

	// 只有当 content 中没有 categories 时才查询
	if _, exists := content["categories"]; !exists {
		db, err := GetProductDB()
		if err != nil {
			hlog.CtxWarnf(ctx, "failed to get product db: %v", err)
			content["categories"] = []Category{}
		} else {
			var categories []Category
			if err := db.Find(&categories).Error; err != nil {
				hlog.CtxWarnf(ctx, "failed to query categories: %v", err)
				content["categories"] = []Category{}
			} else {
				content["categories"] = categories
			}
		}
	}

	return content
}
