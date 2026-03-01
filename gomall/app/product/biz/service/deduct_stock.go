package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
	"github.com/cloudwego/biz-demo/gomall/app/product/infra/mq"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type DeductStockService struct {
	ctx context.Context
}

func NewDeductStockService(ctx context.Context) *DeductStockService {
	return &DeductStockService{ctx: ctx}
}

func (s *DeductStockService) Run(req *product.DeductStockReq) (resp *product.DeductStockResp, err error) {
	resp = &product.DeductStockResp{}

	stockService := NewStockService(s.ctx)

	stockData, err := stockService.GetStock(req.ProductId)
	if err != nil {
		resp.Success = false
		resp.ErrorMessage = "failed to get stock from database: " + err.Error()
		return resp, nil
	}

	if stockData == nil {
		resp.Success = false
		resp.ErrorMessage = "stock record not found in database"
		return resp, nil
	}

	if stockData.Available < int64(req.Quantity) {
		resp.Success = false
		resp.ErrorMessage = "insufficient stock"
		return resp, nil
	}

	remainingStock, redisErr := redis.DeductStock(s.ctx, req.ProductId, int64(req.Quantity))
	
	if redisErr != nil || remainingStock < 0 {
		klog.Infof("Redis stock not ready or insufficient, initializing from database: productId=%d, dbStock=%d", 
			req.ProductId, stockData.Available)
		
		initErr := redis.InitStock(s.ctx, req.ProductId, stockData.Available)
		if initErr != nil {
			resp.Success = false
			resp.ErrorMessage = "failed to init redis stock: " + initErr.Error()
			return resp, nil
		}
		
		remainingStock, redisErr = redis.DeductStock(s.ctx, req.ProductId, int64(req.Quantity))
		if redisErr != nil {
			resp.Success = false
			resp.ErrorMessage = "failed to deduct redis stock: " + redisErr.Error()
			return resp, nil
		}
		
		if remainingStock < 0 {
			resp.Success = false
			resp.ErrorMessage = "insufficient stock after sync"
			return resp, nil
		}
	}

	stockMsg := mq.NewStockMessage(req.ProductId, int64(req.Quantity), "", 0, mq.OperationDeduct)

	err = mq.SendStockDeductMessage(s.ctx, stockMsg)
	if err != nil {
		redis.RestoreStock(s.ctx, req.ProductId, int64(req.Quantity))
		resp.Success = false
		resp.ErrorMessage = "failed to send stock deduct message: " + err.Error()
		return resp, nil
	}

	resp.Success = true
	resp.RemainingStock = remainingStock

	klog.Infof("Stock deduct request sent to RocketMQ: productId=%d, quantity=%d, remaining=%d",
		req.ProductId, req.Quantity, remainingStock)

	return resp, nil
}
