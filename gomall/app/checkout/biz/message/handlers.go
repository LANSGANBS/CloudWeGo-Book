package message

import (
	"context"
	"fmt"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type StockDeductHandler struct{}

func NewStockDeductHandler() *StockDeductHandler {
	return &StockDeductHandler{}
}

func (h *StockDeductHandler) Handle(ctx context.Context, msg *model.LocalMessage) error {
	payload, err := ParsePayload[StockDeductPayload](msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse stock deduct payload: %w", err)
	}
	
	klog.Infof("Processing stock deduct: productId=%d, quantity=%d, orderNo=%s", 
		payload.ProductID, payload.Quantity, payload.OrderNo)
	
	resp, err := rpc.ProductClient.DeductStock(ctx, &product.DeductStockReq{
		ProductId: payload.ProductID,
		Quantity:  payload.Quantity,
	})
	
	if err != nil {
		return fmt.Errorf("failed to deduct stock: %w", err)
	}
	
	if !resp.Success {
		return fmt.Errorf("stock deduct failed: %s", resp.ErrorMessage)
	}
	
	klog.Infof("Stock deducted successfully: productId=%d, remaining=%d", 
		payload.ProductID, resp.RemainingStock)
	
	return nil
}

type StockRestoreHandler struct{}

func NewStockRestoreHandler() *StockRestoreHandler {
	return &StockRestoreHandler{}
}

func (h *StockRestoreHandler) Handle(ctx context.Context, msg *model.LocalMessage) error {
	payload, err := ParsePayload[StockRestorePayload](msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse stock restore payload: %w", err)
	}
	
	klog.Infof("Processing stock restore: productId=%d, quantity=%d, orderNo=%s, reason=%s",
		payload.ProductID, payload.Quantity, payload.OrderNo, payload.Reason)
	
	resp, err := rpc.ProductClient.RestoreStock(ctx, &product.RestoreStockReq{
		ProductId: payload.ProductID,
		Quantity:  payload.Quantity,
	})
	
	if err != nil {
		return fmt.Errorf("failed to restore stock: %w", err)
	}
	
	if !resp.Success {
		return fmt.Errorf("stock restore failed")
	}
	
	klog.Infof("Stock restored successfully: productId=%d", payload.ProductID)
	
	return nil
}

type IncrementSalesHandler struct{}

func NewIncrementSalesHandler() *IncrementSalesHandler {
	return &IncrementSalesHandler{}
}

func (h *IncrementSalesHandler) Handle(ctx context.Context, msg *model.LocalMessage) error {
	payload, err := ParsePayload[struct {
		ProductID uint32 `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}](msg.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse increment sales payload: %w", err)
	}
	
	klog.Infof("Processing increment sales: productId=%d, quantity=%d", 
		payload.ProductID, payload.Quantity)
	
	resp, err := rpc.ProductClient.IncrementSales(ctx, &product.IncrementSalesReq{
		Id:       payload.ProductID,
		Quantity: payload.Quantity,
	})
	
	if err != nil {
		return fmt.Errorf("failed to increment sales: %w", err)
	}
	
	if !resp.Success {
		return fmt.Errorf("increment sales failed")
	}
	
	klog.Infof("Sales incremented successfully: productId=%d", payload.ProductID)
	
	return nil
}
