package service

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
)

type GetProductPriceHistoryService struct {
	ctx context.Context
}

func NewGetProductPriceHistoryService(ctx context.Context) *GetProductPriceHistoryService {
	return &GetProductPriceHistoryService{ctx: ctx}
}

func (s *GetProductPriceHistoryService) Run(req *product.GetProductPriceHistoryReq) (resp *product.GetProductPriceHistoryResp, err error) {
	resp = &product.GetProductPriceHistoryResp{}
	
	if req.ProductId == 0 {
		return nil, kerrors.NewBizStatusError(40000, "product id is required")
	}
	
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	
	histories, err := model.NewProductPriceHistoryQuery(s.ctx, mysql.DB).GetByProductId(req.ProductId, int(limit))
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "failed to get price history")
	}
	
	for _, h := range histories {
		item := &product.PriceHistoryItem{
			Id:           uint32(h.ID),
			ProductId:    h.ProductId,
			ChangeType:   int32(h.ChangeType),
			OperatorName: h.OperatorName,
			Remark:       h.Remark,
			CreatedAt:    h.CreatedAt.Unix(),
		}
		
		if h.OldPrice != nil {
			item.OldPrice = *h.OldPrice
		}
		if h.NewPrice != nil {
			item.NewPrice = *h.NewPrice
		}
		if h.OldDiscountType != nil {
			item.OldDiscountType = int32(*h.OldDiscountType)
		}
		if h.NewDiscountType != nil {
			item.NewDiscountType = int32(*h.NewDiscountType)
		}
		if h.OldDiscountValue != nil {
			item.OldDiscountValue = *h.OldDiscountValue
		}
		if h.NewDiscountValue != nil {
			item.NewDiscountValue = *h.NewDiscountValue
		}
		if h.DiscountStartTime != nil {
			item.DiscountStartTime = h.DiscountStartTime.Unix()
		}
		if h.DiscountEndTime != nil {
			item.DiscountEndTime = h.DiscountEndTime.Unix()
		}
		
		resp.Items = append(resp.Items, item)
	}
	
	return resp, nil
}
