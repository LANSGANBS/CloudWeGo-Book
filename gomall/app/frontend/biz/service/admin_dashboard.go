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
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
)

type AdminDashboardService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAdminDashboardService(Context context.Context, RequestContext *app.RequestContext) *AdminDashboardService {
	return &AdminDashboardService{RequestContext: RequestContext, Context: Context}
}

func (s *AdminDashboardService) Run() (res map[string]any, err error) {
	ctx := s.Context

	// 获取商品总数
	productCount := 0
	products, err := rpc.ProductClient.ListProducts(ctx, &product.ListProductsReq{})
	if err != nil {
		klog.Errorf("获取商品列表失败: %v", err)
	} else {
		productCount = len(products.Products)
	}

	// 获取订单总数和今日销售
	orderCount := 0
	todaySales := 0.0
	orders, err := rpc.OrderClient.ListOrder(ctx, &order.ListOrderReq{})
	if err != nil {
		klog.Errorf("获取订单列表失败: %v", err)
	} else {
		orderCount = len(orders.Orders)
		// 计算今日销售金额
		now := time.Now()
		todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		todayStartTimestamp := int32(todayStart.Unix())
		for _, o := range orders.Orders {
			if o.CreatedAt >= todayStartTimestamp {
				for _, item := range o.OrderItems {
					todaySales += float64(item.Cost)
				}
			}
		}
	}

	// 获取用户总数（只统计普通用户，不包含管理员）
	userCount := 0
	count, err := rpc.GetUserCount(ctx)
	if err != nil {
		klog.Errorf("获取用户总数失败: %v", err)
	} else {
		userCount = int(count)
	}

	return utils.H{
		"title":        "管理后台",
		"productCount": productCount,
		"orderCount":   orderCount,
		"userCount":    userCount,
		"todaySales":   todaySales,
	}, nil
}
