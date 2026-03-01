package message

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type LocalMessageService struct {
	db *gorm.DB
}

func NewLocalMessageService() *LocalMessageService {
	return &LocalMessageService{
		db: mysql.DB,
	}
}

func (s *LocalMessageService) CreateMessage(ctx context.Context, tx *gorm.DB, msg *model.LocalMessage) error {
	if tx == nil {
		tx = s.db
	}
	
	if msg.MessageID == "" {
		msg.MessageID = generateMessageID(msg.MessageType)
	}
	if msg.Status == "" {
		msg.Status = model.MessageStatusPending
	}
	if msg.MaxRetry == 0 {
		msg.MaxRetry = 5
	}
	
	return tx.Create(msg).Error
}

func (s *LocalMessageService) CreateMessageWithBusiness(ctx context.Context, tx *gorm.DB, businessID string, msgType model.MessageType, targetService, targetMethod string, payload interface{}) (*model.LocalMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	msg := &model.LocalMessage{
		MessageID:     generateMessageID(msgType),
		BusinessID:    businessID,
		MessageType:   msgType,
		TargetService: targetService,
		TargetMethod:  targetMethod,
		Payload:       string(payloadBytes),
		Status:        model.MessageStatusPending,
		MaxRetry:      5,
		Priority:      0,
	}
	
	if err := s.CreateMessage(ctx, tx, msg); err != nil {
		return nil, err
	}
	
	return msg, nil
}

func (s *LocalMessageService) UpdateStatus(ctx context.Context, messageID string, status model.MessageStatus, errMsg string) error {
	updates := map[string]interface{}{
		"status":       status,
		"updated_at":   time.Now(),
	}
	
	if errMsg != "" {
		updates["error_message"] = errMsg
	}
	
	if status == model.MessageStatusConfirmed {
		now := time.Now()
		updates["confirm_at"] = &now
	}
	
	return s.db.Model(&model.LocalMessage{}).
		Where("message_id = ?", messageID).
		Updates(updates).Error
}

func (s *LocalMessageService) IncrementRetry(ctx context.Context, messageID string, errMsg string) error {
	now := time.Now()
	nextRetry := calculateNextRetry(now)
	
	return s.db.Model(&model.LocalMessage{}).
		Where("message_id = ?", messageID).
		Updates(map[string]interface{}{
			"retry_count":   gorm.Expr("retry_count + 1"),
			"last_retry_at": &now,
			"next_retry_at": &nextRetry,
			"error_message": errMsg,
			"status":        model.MessageStatusPending,
			"updated_at":    now,
		}).Error
}

func (s *LocalMessageService) GetPendingMessages(ctx context.Context, limit int) ([]model.LocalMessage, error) {
	var messages []model.LocalMessage
	now := time.Now()
	
	err := s.db.Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", model.MessageStatusPending, now).
		Order("priority DESC, created_at ASC").
		Limit(limit).
		Find(&messages).Error
	
	return messages, err
}

func (s *LocalMessageService) GetMessageByID(ctx context.Context, messageID string) (*model.LocalMessage, error) {
	var msg model.LocalMessage
	err := s.db.Where("message_id = ?", messageID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (s *LocalMessageService) GetMessagesByBusinessID(ctx context.Context, businessID string) ([]model.LocalMessage, error) {
	var messages []model.LocalMessage
	err := s.db.Where("business_id = ?", businessID).Order("created_at ASC").Find(&messages).Error
	return messages, err
}

func (s *LocalMessageService) MoveToDLQ(ctx context.Context, messageID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var msg model.LocalMessage
		if err := tx.Where("message_id = ?", messageID).First(&msg).Error; err != nil {
			return err
		}
		
		dlq := &model.MessageDeadLetter{
			OriginalMessageID: msg.MessageID,
			BusinessID:        msg.BusinessID,
			MessageType:       msg.MessageType,
			TargetService:     msg.TargetService,
			TargetMethod:      msg.TargetMethod,
			Payload:           msg.Payload,
			RetryCount:        msg.RetryCount,
			LastError:         msg.ErrorMessage,
			MovedToDLQAt:      time.Now(),
			ManualReview:      true,
		}
		
		if err := tx.Create(dlq).Error; err != nil {
			return err
		}
		
		if err := tx.Delete(&msg).Error; err != nil {
			return err
		}
		
		klog.Infof("Message %s moved to DLQ after %d retries", messageID, msg.RetryCount)
		return nil
	})
}

func (s *LocalMessageService) CancelMessage(ctx context.Context, messageID string, reason string) error {
	return s.db.Model(&model.LocalMessage{}).
		Where("message_id = ? AND status = ?", messageID, model.MessageStatusPending).
		Updates(map[string]interface{}{
			"status":       model.MessageStatusCancelled,
			"error_message": reason,
			"updated_at":   time.Now(),
		}).Error
}

func (s *LocalMessageService) RetryFromDLQ(ctx context.Context, dlqID uint, reviewedBy string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var dlq model.MessageDeadLetter
		if err := tx.First(&dlq, dlqID).Error; err != nil {
			return err
		}
		
		msg := &model.LocalMessage{
			MessageID:     generateMessageID(dlq.MessageType),
			BusinessID:    dlq.BusinessID,
			MessageType:   dlq.MessageType,
			TargetService: dlq.TargetService,
			TargetMethod:  dlq.TargetMethod,
			Payload:       dlq.Payload,
			Status:        model.MessageStatusPending,
			MaxRetry:      5,
			CorrelationID: dlq.OriginalMessageID,
		}
		
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		
		now := time.Now()
		if err := tx.Model(&dlq).Updates(map[string]interface{}{
			"manual_review": false,
			"reviewed_at":   &now,
			"reviewed_by":   reviewedBy,
			"resolution":    "Retried from DLQ",
		}).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func (s *LocalMessageService) LogRetry(ctx context.Context, messageID string, retryCount int, success bool, errMsg string, duration int64) error {
	log := &model.MessageRetryLog{
		MessageID:    messageID,
		RetryCount:   retryCount,
		AttemptAt:    time.Now(),
		Success:      success,
		ErrorMessage: errMsg,
		Duration:     duration,
	}
	return s.db.Create(log).Error
}

func (s *LocalMessageService) GetDLQMessages(ctx context.Context, limit int, offset int) ([]model.MessageDeadLetter, int64, error) {
	var messages []model.MessageDeadLetter
	var total int64
	
	s.db.Model(&model.MessageDeadLetter{}).Count(&total)
	
	err := s.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	
	return messages, total, err
}

func (s *LocalMessageService) GetMessageStats(ctx context.Context) (*MessageStats, error) {
	stats := &MessageStats{}
	
	rows, err := s.db.Model(&model.LocalMessage{}).
		Select("status, count(*) as count").
		Group("status").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch model.MessageStatus(status) {
		case model.MessageStatusPending:
			stats.Pending = count
		case model.MessageStatusSent:
			stats.Sent = count
		case model.MessageStatusConfirmed:
			stats.Confirmed = count
		case model.MessageStatusFailed:
			stats.Failed = count
		case model.MessageStatusCancelled:
			stats.Cancelled = count
		}
	}
	
	s.db.Model(&model.MessageDeadLetter{}).Count(&stats.DLQCount)
	
	return stats, nil
}

type MessageStats struct {
	Pending   int64 `json:"pending"`
	Sent      int64 `json:"sent"`
	Confirmed int64 `json:"confirmed"`
	Failed    int64 `json:"failed"`
	Cancelled int64 `json:"cancelled"`
	DLQCount  int64 `json:"dlq_count"`
}

func generateMessageID(msgType model.MessageType) string {
	return fmt.Sprintf("msg_%s_%d_%d", msgType, time.Now().UnixNano(), rand.Intn(10000))
}

func calculateNextRetry(now time.Time) time.Time {
	return now.Add(time.Minute * 2)
}
