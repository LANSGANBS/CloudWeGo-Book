package service

import (
	"context"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"gorm.io/gorm"
)

type CancelDiscountService struct {
	ctx context.Context
}

func NewCancelDiscountService(ctx context.Context) *CancelDiscountService {
	return &CancelDiscountService{ctx: ctx}
}

func (s *CancelDiscountService) Run(req *product.CancelDiscountReq) (resp *product.CancelDiscountResp, err error) {
	resp = &product.CancelDiscountResp{}
	
	if req.ProductId == 0 {
		return nil, kerrors.NewBizStatusError(40000, "product id is required")
	}
	
	p, err := model.NewProductQuery(s.ctx, mysql.DB).GetById(int(req.ProductId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(40400, "product not found")
	}
	
	if p.DiscountType == model.DiscountTypeNone {
		resp.Success = true
		resp.Message = "商品当前无折扣"
		return resp, nil
	}
	
	oldDiscountType := p.DiscountType
	oldDiscountValue := p.DiscountValue
	oldPrice := p.Price
	
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"discount_type":       model.DiscountTypeNone,
			"discount_value":      0,
			"discount_start_time": nil,
			"discount_end_time":   nil,
			"original_price":      nil,
			"updated_at":          time.Now(),
		}
		
		if err := tx.Model(&p).Updates(updates).Error; err != nil {
			return err
		}
		
		odt := oldDiscountType
		odv := oldDiscountValue
		ndt := model.DiscountTypeNone
		var ndv float32 = 0
		
		return model.RecordPriceChange(tx, s.ctx, req.ProductId, model.PriceChangeTypeCancelDisc, &oldPrice, &p.Price, &odt, &ndt, &odv, &ndv, nil, nil, 0, "admin", "取消折扣")
	})
	
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "failed to cancel discount: "+err.Error())
	}
	
	resp.Success = true
	resp.Message = "折扣已取消"
	return resp, nil
}
