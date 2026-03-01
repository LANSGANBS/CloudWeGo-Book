package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusConfirmed MessageStatus = "confirmed"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusCancelled MessageStatus = "cancelled"
)

type MessageType string

const (
	MessageTypeStockDeduct  MessageType = "stock_deduct"
	MessageTypeStockRestore MessageType = "stock_restore"
	MessageTypeOrderCreate  MessageType = "order_create"
	MessageTypePayment      MessageType = "payment"
	MessageTypeEmail        MessageType = "email"
)

type LocalMessage struct {
	Base
	MessageID      string        `gorm:"uniqueIndex;size:128;not null" json:"message_id"`
	BusinessID     string        `gorm:"index;size:128;not null" json:"business_id"`
	MessageType    MessageType   `gorm:"size:32;not null" json:"message_type"`
	TargetService  string        `gorm:"size:64;not null" json:"target_service"`
	TargetMethod   string        `gorm:"size:64;not null" json:"target_method"`
	Payload        string        `gorm:"type:text;not null" json:"payload"`
	Status         MessageStatus `gorm:"size:32;default:'pending';index" json:"status"`
	RetryCount     int           `gorm:"default:0" json:"retry_count"`
	MaxRetry       int           `gorm:"default:5" json:"max_retry"`
	NextRetryAt    *time.Time    `json:"next_retry_at"`
	LastRetryAt    *time.Time    `json:"last_retry_at"`
	ErrorMessage   string        `gorm:"type:text" json:"error_message"`
	ConfirmAt      *time.Time    `json:"confirm_at"`
	Priority       int           `gorm:"default:0;index" json:"priority"`
	CorrelationID  string        `gorm:"size:128;index" json:"correlation_id"`
	CallbackURL    string        `gorm:"size:256" json:"callback_url"`
	CallbackStatus string        `gorm:"size:32" json:"callback_status"`
}

func (LocalMessage) TableName() string {
	return "local_message"
}

type MessageDeadLetter struct {
	Base
	OriginalMessageID string      `gorm:"index;size:128;not null" json:"original_message_id"`
	BusinessID        string      `gorm:"index;size:128;not null" json:"business_id"`
	MessageType       MessageType `gorm:"size:32;not null" json:"message_type"`
	TargetService     string      `gorm:"size:64;not null" json:"target_service"`
	TargetMethod      string      `gorm:"size:64;not null" json:"target_method"`
	Payload           string      `gorm:"type:text;not null" json:"payload"`
	RetryCount        int         `json:"retry_count"`
	LastError         string      `gorm:"type:text" json:"last_error"`
	MovedToDLQAt      time.Time   `json:"moved_to_dlq_at"`
	ManualReview      bool        `gorm:"default:false" json:"manual_review"`
	ReviewedAt        *time.Time  `json:"reviewed_at"`
	ReviewedBy        string      `gorm:"size:64" json:"reviewed_by"`
	Resolution        string      `gorm:"type:text" json:"resolution"`
}

func (MessageDeadLetter) TableName() string {
	return "message_dead_letter"
}

type MessageRetryLog struct {
	Base
	MessageID    string    `gorm:"index;size:128;not null" json:"message_id"`
	RetryCount   int       `json:"retry_count"`
	AttemptAt    time.Time `json:"attempt_at"`
	Success      bool      `json:"success"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	Duration     int64     `json:"duration"`
}

func (MessageRetryLog) TableName() string {
	return "message_retry_log"
}
