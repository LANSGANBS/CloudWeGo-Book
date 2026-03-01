package message

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
)

type MessageMonitor struct {
	service *LocalMessageService
}

func NewMessageMonitor(service *LocalMessageService) *MessageMonitor {
	return &MessageMonitor{
		service: service,
	}
}

type MessageDetail struct {
	Message         *model.LocalMessage     `json:"message"`
	RetryLogs       []model.MessageRetryLog `json:"retry_logs"`
	CanRetry        bool                    `json:"can_retry"`
	CanCancel       bool                    `json:"can_cancel"`
	TimeInStatus    time.Duration           `json:"time_in_status"`
	EstimatedRetry  *time.Time              `json:"estimated_retry,omitempty"`
}

func (m *MessageMonitor) GetMessageDetail(ctx context.Context, messageID string) (*MessageDetail, error) {
	msg, err := m.service.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	
	var retryLogs []model.MessageRetryLog
	m.service.db.Where("message_id = ?", messageID).Order("attempt_at DESC").Find(&retryLogs)
	
	detail := &MessageDetail{
		Message:   msg,
		RetryLogs: retryLogs,
		CanRetry:  msg.Status == model.MessageStatusPending || msg.Status == model.MessageStatusFailed,
		CanCancel: msg.Status == model.MessageStatusPending,
	}
	
	if msg.Status == model.MessageStatusPending {
		detail.TimeInStatus = time.Since(msg.UpdatedAt)
		if msg.NextRetryAt != nil {
			detail.EstimatedRetry = msg.NextRetryAt
		}
	}
	
	return detail, nil
}

type MessageListFilter struct {
	Status      *model.MessageStatus `form:"status"`
	MessageType *model.MessageType   `form:"message_type"`
	BusinessID  string               `form:"business_id"`
	StartDate   *time.Time           `form:"start_date"`
	EndDate     *time.Time           `form:"end_date"`
	Page        int                  `form:"page"`
	PageSize    int                  `form:"page_size"`
}

type MessageListResult struct {
	Messages  []model.LocalMessage `json:"messages"`
	Total     int64                `json:"total"`
	Page      int                  `json:"page"`
	PageSize  int                  `json:"page_size"`
	TotalPage int                  `json:"total_page"`
}

func (m *MessageMonitor) ListMessages(ctx context.Context, filter *MessageListFilter) (*MessageListResult, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	
	query := m.service.db.Model(&model.LocalMessage{})
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.MessageType != nil {
		query = query.Where("message_type = ?", *filter.MessageType)
	}
	if filter.BusinessID != "" {
		query = query.Where("business_id = ?", filter.BusinessID)
	}
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}
	
	var total int64
	query.Count(&total)
	
	var messages []model.LocalMessage
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&messages).Error; err != nil {
		return nil, err
	}
	
	totalPage := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPage++
	}
	
	return &MessageListResult{
		Messages:  messages,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPage,
	}, nil
}

type DashboardStats struct {
	TotalMessages   int64                    `json:"total_messages"`
	PendingMessages int64                    `json:"pending_messages"`
	ConfirmedMessages int64                  `json:"confirmed_messages"`
	FailedMessages  int64                    `json:"failed_messages"`
	DLQCount        int64                    `json:"dlq_count"`
	RecentFailures  []model.LocalMessage     `json:"recent_failures"`
	StatusTrend     []StatusTrendItem        `json:"status_trend"`
	MessageTypeDist []MessageTypeDistItem    `json:"message_type_distribution"`
}

type StatusTrendItem struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type MessageTypeDistItem struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func (m *MessageMonitor) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats, err := m.service.GetMessageStats(ctx)
	if err != nil {
		return nil, err
	}
	
	var recentFailures []model.LocalMessage
	m.service.db.Where("status = ?", model.MessageStatusFailed).
		Order("updated_at DESC").
		Limit(10).
		Find(&recentFailures)
	
	var typeDistribution []MessageTypeDistItem
	rows, err := m.service.db.Model(&model.LocalMessage{}).
		Select("message_type, count(*) as count").
		Group("message_type").
		Rows()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var msgType string
			var count int64
			if err := rows.Scan(&msgType, &count); err == nil {
				typeDistribution = append(typeDistribution, MessageTypeDistItem{
					Type:  msgType,
					Count: count,
				})
			}
		}
	}
	
	var statusTrend []StatusTrendItem
	trendRows, err := m.service.db.Model(&model.LocalMessage{}).
		Select("DATE(created_at) as date, count(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("date ASC").
		Rows()
	if err == nil {
		defer trendRows.Close()
		for trendRows.Next() {
			var date string
			var count int64
			if err := trendRows.Scan(&date, &count); err == nil {
				statusTrend = append(statusTrend, StatusTrendItem{
					Date:  date,
					Count: count,
				})
			}
		}
	}
	
	return &DashboardStats{
		TotalMessages:     stats.Pending + stats.Sent + stats.Confirmed + stats.Failed + stats.Cancelled,
		PendingMessages:   stats.Pending,
		ConfirmedMessages: stats.Confirmed,
		FailedMessages:    stats.Failed,
		DLQCount:          stats.DLQCount,
		RecentFailures:    recentFailures,
		StatusTrend:       statusTrend,
		MessageTypeDist:   typeDistribution,
	}, nil
}

type DLQDetail struct {
	Message      *model.MessageDeadLetter `json:"message"`
	CanRetry     bool                     `json:"can_retry"`
	TimeInDLQ    time.Duration            `json:"time_in_dlq"`
}

func (m *MessageMonitor) GetDLQDetail(ctx context.Context, id uint) (*DLQDetail, error) {
	var dlq model.MessageDeadLetter
	if err := m.service.db.First(&dlq, id).Error; err != nil {
		return nil, err
	}
	
	return &DLQDetail{
		Message:   &dlq,
		CanRetry:  true,
		TimeInDLQ: time.Since(dlq.MovedToDLQAt),
	}, nil
}

type DLQListFilter struct {
	ManualReview *bool  `form:"manual_review"`
	MessageType  string `form:"message_type"`
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
}

type DLQListResult struct {
	Messages  []model.MessageDeadLetter `json:"messages"`
	Total     int64                     `json:"total"`
	Page      int                       `json:"page"`
	PageSize  int                       `json:"page_size"`
	TotalPage int                       `json:"total_page"`
}

func (m *MessageMonitor) ListDLQ(ctx context.Context, filter *DLQListFilter) (*DLQListResult, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	
	query := m.service.db.Model(&model.MessageDeadLetter{})
	
	if filter.ManualReview != nil {
		query = query.Where("manual_review = ?", *filter.ManualReview)
	}
	if filter.MessageType != "" {
		query = query.Where("message_type = ?", filter.MessageType)
	}
	
	var total int64
	query.Count(&total)
	
	var messages []model.MessageDeadLetter
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("moved_to_dlq_at DESC").Offset(offset).Limit(filter.PageSize).Find(&messages).Error; err != nil {
		return nil, err
	}
	
	totalPage := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPage++
	}
	
	return &DLQListResult{
		Messages:  messages,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPage,
	}, nil
}

func (m *MessageMonitor) RetryDLQMessage(ctx context.Context, id uint, reviewedBy string) error {
	return m.service.RetryFromDLQ(ctx, id, reviewedBy)
}

func (m *MessageMonitor) DeleteDLQMessage(ctx context.Context, id uint, reason string) error {
	return m.service.db.Model(&model.MessageDeadLetter{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"resolution":  reason,
			"deleted_at":  time.Now(),
		}).Error
}

func (m *MessageMonitor) BatchRetryMessages(ctx context.Context, messageIDs []string) (int, []string) {
	successCount := 0
	var failedIDs []string
	
	for _, id := range messageIDs {
		if err := m.service.UpdateStatus(ctx, id, model.MessageStatusPending, ""); err != nil {
			klog.Errorf("Failed to reset message %s for retry: %v", id, err)
			failedIDs = append(failedIDs, id)
		} else {
			successCount++
		}
	}
	
	return successCount, failedIDs
}

func (m *MessageMonitor) BatchCancelMessages(ctx context.Context, messageIDs []string, reason string) (int, []string) {
	successCount := 0
	var failedIDs []string
	
	for _, id := range messageIDs {
		if err := m.service.CancelMessage(ctx, id, reason); err != nil {
			klog.Errorf("Failed to cancel message %s: %v", id, err)
			failedIDs = append(failedIDs, id)
		} else {
			successCount++
		}
	}
	
	return successCount, failedIDs
}

func (m *MessageMonitor) ExportMessages(ctx context.Context, filter *MessageListFilter) ([]byte, error) {
	filter.PageSize = 1000
	result, err := m.ListMessages(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	return json.MarshalIndent(result.Messages, "", "  ")
}

func (m *MessageMonitor) GetHealthStatus(ctx context.Context) (*HealthStatus, error) {
	stats, err := m.service.GetMessageStats(ctx)
	if err != nil {
		return nil, err
	}
	
	health := &HealthStatus{
		Status:      "healthy",
		PendingRate: float64(stats.Pending) / float64(stats.Pending+stats.Confirmed+stats.Failed+1) * 100,
		FailureRate: float64(stats.Failed) / float64(stats.Pending+stats.Confirmed+stats.Failed+1) * 100,
	}
	
	if stats.Pending > 1000 {
		health.Status = "warning"
		health.Warnings = append(health.Warnings, "High number of pending messages")
	}
	
	if stats.Failed > 100 {
		health.Status = "warning"
		health.Warnings = append(health.Warnings, "High number of failed messages")
	}
	
	if stats.DLQCount > 50 {
		health.Status = "critical"
		health.Warnings = append(health.Warnings, "High number of messages in DLQ")
	}
	
	return health, nil
}

type HealthStatus struct {
	Status      string   `json:"status"`
	PendingRate float64  `json:"pending_rate"`
	FailureRate float64  `json:"failure_rate"`
	Warnings    []string `json:"warnings,omitempty"`
}
