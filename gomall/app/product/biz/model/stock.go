package model

import (
	"time"
)

type Stock struct {
	Base
	ProductId    uint32  `gorm:"uniqueIndex;not null" json:"product_id"`
	Quantity     int64   `gorm:"default:0" json:"quantity"`
	Reserved     int64   `gorm:"default:0" json:"reserved"`
	Available    int64   `gorm:"default:0" json:"available"`
	MinStock     int64   `gorm:"default:10" json:"min_stock"`
	MaxStock     int64   `gorm:"default:1000" json:"max_stock"`
	SafetyStock  int64   `gorm:"default:20" json:"safety_stock"`
	Unit         string  `gorm:"size:20;default:'件'" json:"unit"`
	WarehouseId  uint32  `gorm:"default:1" json:"warehouse_id"`
	Location     string  `gorm:"size:100" json:"location"`
	BatchNo      string  `gorm:"size:50" json:"batch_no"`
	ExpiredAt    *time.Time `json:"expired_at"`
	Status       int8    `gorm:"default:1" json:"status"`
}

func (Stock) TableName() string {
	return "stock"
}

type StockLog struct {
	Base
	ProductId    uint32    `gorm:"index;not null" json:"product_id"`
	OrderNo      string    `gorm:"size:64;index" json:"order_no"`
	ChangeType   int8      `gorm:"not null" json:"change_type"`
	ChangeQty    int64     `gorm:"not null" json:"change_qty"`
	BeforeQty    int64     `gorm:"not null" json:"before_qty"`
	AfterQty     int64     `gorm:"not null" json:"after_qty"`
	OperatorId   uint32    `json:"operator_id"`
	OperatorName string    `gorm:"size:50" json:"operator_name"`
	Remark       string    `gorm:"size:255" json:"remark"`
	WarehouseId  uint32    `gorm:"default:1" json:"warehouse_id"`
}

func (StockLog) TableName() string {
	return "stock_log"
}

type StockAlert struct {
	Base
	ProductId    uint32    `gorm:"index;not null" json:"product_id"`
	AlertType    int8      `gorm:"not null" json:"alert_type"`
	AlertLevel   int8      `gorm:"default:1" json:"alert_level"`
	Threshold    int64     `gorm:"not null" json:"threshold"`
	CurrentValue int64     `gorm:"not null" json:"current_value"`
	Status       int8      `gorm:"default:0" json:"status"`
	HandledAt    *time.Time `json:"handled_at"`
	HandlerId    uint32    `json:"handler_id"`
	HandlerName  string    `gorm:"size:50" json:"handler_name"`
	Remark       string    `gorm:"size:255" json:"remark"`
}

func (StockAlert) TableName() string {
	return "stock_alert"
}

type StockCheck struct {
	Base
	CheckNo      string    `gorm:"size:64;uniqueIndex" json:"check_no"`
	WarehouseId  uint32    `gorm:"default:1" json:"warehouse_id"`
	Status       int8      `gorm:"default:0" json:"status"`
	TotalItems   int       `gorm:"default:0" json:"total_items"`
	DiffItems    int       `gorm:"default:0" json:"diff_items"`
	OperatorId   uint32    `json:"operator_id"`
	OperatorName string    `gorm:"size:50" json:"operator_name"`
	Remark       string    `gorm:"size:255" json:"remark"`
	FinishedAt   *time.Time `json:"finished_at"`
}

func (StockCheck) TableName() string {
	return "stock_check"
}

type StockCheckItem struct {
	Base
	CheckId      uint32    `gorm:"index;not null" json:"check_id"`
	ProductId    uint32    `gorm:"index;not null" json:"product_id"`
	SystemQty    int64     `gorm:"not null" json:"system_qty"`
	ActualQty    int64     `gorm:"not null" json:"actual_qty"`
	DiffQty      int64     `gorm:"not null" json:"diff_qty"`
	Remark       string    `gorm:"size:255" json:"remark"`
}

func (StockCheckItem) TableName() string {
	return "stock_check_item"
}

const (
	StockStatusNormal   = 1
	StockStatusLocked   = 2
	StockStatusDisabled = 0
)

const (
	ChangeTypePurchase     = 1
	ChangeTypeSale         = 2
	ChangeTypeReturn       = 3
	ChangeTypeAdjust       = 4
	ChangeTypeCheck        = 5
	ChangeTypeReserve      = 6
	ChangeTypeRelease      = 7
	ChangeTypeDamage       = 8
	ChangeTypeTransfer     = 9
)

const (
	AlertTypeLowStock    = 1
	AlertTypeOverStock   = 2
	AlertTypeExpiring    = 3
	AlertTypeExpired     = 4
)

const (
	AlertLevelInfo    = 1
	AlertLevelWarning = 2
	AlertLevelDanger  = 3
)

const (
	AlertStatusPending  = 0
	AlertStatusHandled  = 1
	AlertStatusIgnored  = 2
)

const (
	CheckStatusPending   = 0
	CheckStatusProgress  = 1
	CheckStatusFinished  = 2
	CheckStatusCancelled = 3
)

type StockMessageLog struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	MessageId   string    `gorm:"uniqueIndex;size:128;not null" json:"message_id"`
	ProductId   uint32    `gorm:"index" json:"product_id"`
	Operation   string    `gorm:"size:32" json:"operation"`
	Quantity    int64     `json:"quantity"`
	Status      string    `gorm:"size:32;default:'processed'" json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (StockMessageLog) TableName() string {
	return "stock_message_log"
}

type StockDLQ struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	MessageId    string    `gorm:"uniqueIndex;size:128" json:"message_id"`
	ProductId    uint32    `gorm:"index" json:"product_id"`
	Quantity     int64     `json:"quantity"`
	OrderNo      string    `gorm:"size:64" json:"order_no"`
	UserId       uint32    `json:"user_id"`
	Operation    string    `gorm:"size:32" json:"operation"`
	RetryCount   int       `json:"retry_count"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (StockDLQ) TableName() string {
	return "stock_dlq"
}
