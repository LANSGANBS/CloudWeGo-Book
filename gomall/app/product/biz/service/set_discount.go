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

type SetDiscountService struct {
	ctx context.Context
}

func NewSetDiscountService(ctx context.Context) *SetDiscountService {
	return &SetDiscountService{ctx: ctx}
}

func (s *SetDiscountService) Run(req *product.SetDiscountReq) (resp *product.SetDiscountResp, err error) {
	resp = &product.SetDiscountResp{}
	
	if req.ProductId == 0 {
		return nil, kerrors.NewBizStatusError(40000, "product id is required")
	}
	
	if req.DiscountType < 0 || req.DiscountType > 2 {
		return nil, kerrors.NewBizStatusError(40001, "invalid discount type")
	}
	
	if req.DiscountType != 0 && req.DiscountValue <= 0 {
		return nil, kerrors.NewBizStatusError(40002, "discount value must be greater than 0")
	}
	
	if req.DiscountType == 1 && req.DiscountValue >= 1 {
		return nil, kerrors.NewBizStatusError(40003, "discount rate must be less than 1 (e.g., 0.8 for 20% off)")
	}
	
	if req.StartTime > 0 && req.EndTime > 0 && req.StartTime >= req.EndTime {
		return nil, kerrors.NewBizStatusError(40004, "start time must be before end time")
	}
	
	p, err := model.NewProductQuery(s.ctx, mysql.DB).GetById(int(req.ProductId))
	if err != nil {
		return nil, kerrors.NewBizStatusError(40400, "product not found")
	}
	
	oldDiscountType := p.DiscountType
	oldDiscountValue := p.DiscountValue
	oldPrice := p.Price
	
	var startTime, endTime *time.Time
	if req.StartTime > 0 {
		t := time.Unix(req.StartTime, 0)
		startTime = &t
	}
	if req.EndTime > 0 {
		t := time.Unix(req.EndTime, 0)
		endTime = &t
	}
	
	var originalPrice *float32
	if req.DiscountType != 0 && p.OriginalPrice == nil {
		op := p.Price
		originalPrice = &op
	} else if req.DiscountType == 0 {
		originalPrice = nil
	} else {
		originalPrice = p.OriginalPrice
	}
	
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"discount_type":       req.DiscountType,
			"discount_value":      req.DiscountValue,
			"discount_start_time": startTime,
			"discount_end_time":   endTime,
			"original_price":      originalPrice,
			"updated_at":          time.Now(),
		}
		
		if err := tx.Model(&p).Updates(updates).Error; err != nil {
			return err
		}
		
		var changeType int8
		if startTime != nil && endTime != nil {
			changeType = model.PriceChangeTypeSetFlashSale
		} else if req.DiscountType != 0 {
			changeType = model.PriceChangeTypeSetDiscount
		} else {
			changeType = model.PriceChangeTypeCancelDisc
		}
		
		odt := oldDiscountType
		odv := oldDiscountValue
		ndt := int8(req.DiscountType)
		ndv := req.DiscountValue
		
		return model.RecordPriceChange(tx, s.ctx, req.ProductId, changeType, &oldPrice, &p.Price, &odt, &ndt, &odv, &ndv, startTime, endTime, 0, "admin", "")
	})
	
	if err != nil {
		return nil, kerrors.NewBizStatusError(50000, "failed to set discount: "+err.Error())
	}
	
	resp.Success = true
	resp.Message = "折扣设置成功"
	return resp, nil
}
