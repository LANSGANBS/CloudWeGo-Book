package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
	"github.com/cloudwego/biz-demo/gomall/app/product/infra/mq"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RestoreStockService struct {
	ctx context.Context
}

func NewRestoreStockService(ctx context.Context) *RestoreStockService {
	return &RestoreStockService{ctx: ctx}
}

func (s *RestoreStockService) Run(req *product.RestoreStockReq) (resp *product.RestoreStockResp, err error) {
	resp = &product.RestoreStockResp{}

	stockService := NewStockService(s.ctx)

	stockData, err := stockService.GetStock(req.ProductId)
	if err != nil {
		klog.Warnf("Failed to get stock from database: %v", err)
	}

	err = redis.RestoreStock(s.ctx, req.ProductId, int64(req.Quantity))
	if err != nil {
		klog.Warnf("Redis restore stock failed: %v", err)
		if stockData != nil {
			redis.InitStock(s.ctx, req.ProductId, stockData.Available+int64(req.Quantity))
		}
	}

	stockMsg := mq.NewStockMessage(req.ProductId, int64(req.Quantity), "", 0, mq.OperationRestore)

	err = mq.SendStockRestoreMessage(s.ctx, stockMsg)
	if err != nil {
		resp.Success = false
		return resp, nil
	}

	resp.Success = true

	klog.Infof("Stock restore request sent to RocketMQ: productId=%d, quantity=%d", req.ProductId, req.Quantity)

	return resp, nil
}
