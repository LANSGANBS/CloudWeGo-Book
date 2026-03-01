package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
)

type GetStockService struct {
	ctx context.Context
}

func NewGetStockService(ctx context.Context) *GetStockService {
	return &GetStockService{ctx: ctx}
}

func (s *GetStockService) Run(req *product.GetStockReq) (resp *product.GetStockResp, err error) {
	resp = &product.GetStockResp{}
	
	stock, err := redis.GetStock(s.ctx, req.ProductId)
	if err != nil {
		resp.Exists = false
		resp.Stock = 0
		return resp, nil
	}
	
	resp.Exists = true
	resp.Stock = stock
	
	return resp, nil
}
