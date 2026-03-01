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

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type IncrementSalesService struct {
	ctx context.Context
}

func NewIncrementSalesService(ctx context.Context) *IncrementSalesService {
	return &IncrementSalesService{ctx: ctx}
}

func (s *IncrementSalesService) Run(req *product.IncrementSalesReq) (resp *product.IncrementSalesResp, err error) {
	if req.Id == 0 || req.Quantity <= 0 {
		klog.Errorf("invalid request: id=%d, quantity=%d", req.Id, req.Quantity)
		return &product.IncrementSalesResp{Success: false}, nil
	}

	err = model.IncrementSales(mysql.DB, s.ctx, int(req.Id), int(req.Quantity))
	if err != nil {
		klog.Errorf("failed to increment sales: %v", err)
		return &product.IncrementSalesResp{Success: false}, err
	}

	return &product.IncrementSalesResp{Success: true}, nil
}
