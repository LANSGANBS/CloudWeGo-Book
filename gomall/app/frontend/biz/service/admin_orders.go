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

	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
)

type AdminOrdersService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAdminOrdersService(Context context.Context, RequestContext *app.RequestContext) *AdminOrdersService {
	return &AdminOrdersService{RequestContext: RequestContext, Context: Context}
}

func (s *AdminOrdersService) Run() (res map[string]any, err error) {
	ctx := s.Context

	// 获取订单列表
	orders, err := rpc.OrderClient.ListOrder(ctx, &order.ListOrderReq{})
	if err != nil {
		klog.Errorf("获取订单列表失败: %v", err)
		// 服务不可用时返回空列表，不返回错误
		return utils.H{
			"title":  "订单管理",
			"orders": []*order.Order{},
		}, nil
	}

	return utils.H{
		"title":  "订单管理",
		"orders": orders.Orders,
	}, nil
}
