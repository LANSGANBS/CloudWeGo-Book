package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	PriceChangeTypeNormalAdjust int8 = 1
	PriceChangeTypeSetDiscount  int8 = 2
	PriceChangeTypeSetFlashSale int8 = 3
	PriceChangeTypeCancelDisc   int8 = 4
	PriceChangeTypeFlashSaleExp int8 = 5
)

type ProductPriceHistory struct {
	ID                uint32         `gorm:"primarykey;column:id"`
	ProductId         uint32         `json:"product_id" gorm:"column:product_id;not null"`
	ChangeType        int8           `json:"change_type" gorm:"column:change_type;not null"`
	OldPrice          *float32       `json:"old_price" gorm:"column:old_price"`
	NewPrice          *float32       `json:"new_price" gorm:"column:new_price"`
	OldDiscountType   *int8          `json:"old_discount_type" gorm:"column:old_discount_type"`
	NewDiscountType   *int8          `json:"new_discount_type" gorm:"column:new_discount_type"`
	OldDiscountValue  *float32       `json:"old_discount_value" gorm:"column:old_discount_value"`
	NewDiscountValue  *float32       `json:"new_discount_value" gorm:"column:new_discount_value"`
	DiscountStartTime *time.Time     `json:"discount_start_time" gorm:"column:discount_start_time"`
	DiscountEndTime   *time.Time     `json:"discount_end_time" gorm:"column:discount_end_time"`
	OperatorId        uint32         `json:"operator_id" gorm:"column:operator_id;default:0"`
	OperatorName      string         `json:"operator_name" gorm:"column:operator_name;default:''"`
	Remark            string         `json:"remark" gorm:"column:remark;default:''"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ProductPriceHistory) TableName() string {
	return "product_price_history"
}

type ProductPriceHistoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProductPriceHistoryQuery(ctx context.Context, db *gorm.DB) ProductPriceHistoryQuery {
	return ProductPriceHistoryQuery{ctx: ctx, db: db}
}

func (q ProductPriceHistoryQuery) Create(history *ProductPriceHistory) error {
	return q.db.WithContext(q.ctx).Create(history).Error
}

func (q ProductPriceHistoryQuery) GetByProductId(productId uint32, limit int) ([]ProductPriceHistory, error) {
	var histories []ProductPriceHistory
	err := q.db.WithContext(q.ctx).
		Where("product_id = ?", productId).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (q ProductPriceHistoryQuery) GetByProductIdAtTime(productId uint32, atTime time.Time) (*ProductPriceHistory, error) {
	var history ProductPriceHistory
	err := q.db.WithContext(q.ctx).
		Where("product_id = ?", productId).
		Where("created_at <= ?", atTime).
		Order("created_at DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func RecordPriceChange(db *gorm.DB, ctx context.Context, productId uint32, changeType int8, oldPrice, newPrice *float32, oldDiscountType, newDiscountType *int8, oldDiscountValue, newDiscountValue *float32, startTime, endTime *time.Time, operatorId uint32, operatorName, remark string) error {
	history := &ProductPriceHistory{
		ProductId:         productId,
		ChangeType:        changeType,
		OldPrice:          oldPrice,
		NewPrice:          newPrice,
		OldDiscountType:   oldDiscountType,
		NewDiscountType:   newDiscountType,
		OldDiscountValue:  oldDiscountValue,
		NewDiscountValue:  newDiscountValue,
		DiscountStartTime: startTime,
		DiscountEndTime:   endTime,
		OperatorId:        operatorId,
		OperatorName:      operatorName,
		Remark:            remark,
	}
	return NewProductPriceHistoryQuery(ctx, db).Create(history)
}
