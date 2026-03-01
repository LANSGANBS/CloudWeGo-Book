package message

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/stretchr/testify/assert"
)

func TestMessageMonitor_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service := GetMessageService()
	if service == nil {
		t.Skip("Message service not initialized")
	}
	
	monitor := NewMessageMonitor(service)
	ctx := context.Background()
	
	t.Run("GetMessageDetail", func(t *testing.T) {
		msg := &model.LocalMessage{
			MessageID:     "monitor-test-001",
			BusinessID:    "order-monitor-001",
			MessageType:   model.MessageTypeStockDeduct,
			TargetService: "product",
			TargetMethod:  "DeductStock",
			Payload:       `{"product_id": 1}`,
			Status:        model.MessageStatusPending,
		}
		service.CreateMessage(ctx, nil, msg)
		
		detail, err := monitor.GetMessageDetail(ctx, msg.MessageID)
		assert.NoError(t, err)
		assert.NotNil(t, detail)
		assert.Equal(t, msg.MessageID, detail.Message.MessageID)
		assert.True(t, detail.CanRetry)
		assert.True(t, detail.CanCancel)
	})
	
	t.Run("ListMessages", func(t *testing.T) {
		filter := &MessageListFilter{
			Page:     1,
			PageSize: 10,
		}
		
		result, err := monitor.ListMessages(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, int64(0))
	})
	
	t.Run("GetDashboardStats", func(t *testing.T) {
		stats, err := monitor.GetDashboardStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
	})
	
	t.Run("GetHealthStatus", func(t *testing.T) {
		health, err := monitor.GetHealthStatus(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, health)
		assert.NotEmpty(t, health.Status)
	})
}

func TestMessageMonitor_ListMessagesWithFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service := GetMessageService()
	if service == nil {
		t.Skip("Message service not initialized")
	}
	
	monitor := NewMessageMonitor(service)
	ctx := context.Background()
	
	pending := model.MessageStatusPending
	filter := &MessageListFilter{
		Status:   &pending,
		Page:     1,
		PageSize: 10,
	}
	
	result, err := monitor.ListMessages(ctx, filter)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestMessageMonitor_DLQOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service := GetMessageService()
	if service == nil {
		t.Skip("Message service not initialized")
	}
	
	monitor := NewMessageMonitor(service)
	ctx := context.Background()
	
	t.Run("ListDLQ", func(t *testing.T) {
		filter := &DLQListFilter{
			Page:     1,
			PageSize: 10,
		}
		
		result, err := monitor.ListDLQ(ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
	
	t.Run("RetryFromDLQ", func(t *testing.T) {
		dlq := &model.MessageDeadLetter{
			OriginalMessageID: "dlq-retry-test-001",
			BusinessID:        "order-dlq-001",
			MessageType:       model.MessageTypeStockDeduct,
			TargetService:     "product",
			TargetMethod:      "DeductStock",
			Payload:           `{"product_id": 1}`,
			RetryCount:        5,
			LastError:         "timeout",
			MovedToDLQAt:      time.Now(),
			ManualReview:      true,
		}
		service.db.Create(dlq)
		
		err := monitor.RetryDLQMessage(ctx, dlq.ID, "admin")
		assert.NoError(t, err)
		
		var newMsg model.LocalMessage
		err = service.db.Where("correlation_id = ?", "dlq-retry-test-001").First(&newMsg).Error
		assert.NoError(t, err)
		assert.Equal(t, model.MessageStatusPending, newMsg.Status)
	})
}

func TestMessageMonitor_BatchOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service := GetMessageService()
	if service == nil {
		t.Skip("Message service not initialized")
	}
	
	monitor := NewMessageMonitor(service)
	ctx := context.Background()
	
	var messageIDs []string
	for i := 0; i < 3; i++ {
		msg := &model.LocalMessage{
			MessageID:     string(rune('x' + i)),
			BusinessID:    "order-batch-001",
			MessageType:   model.MessageTypeStockDeduct,
			TargetService: "product",
			TargetMethod:  "DeductStock",
			Payload:       `{}`,
			Status:        model.MessageStatusFailed,
		}
		service.CreateMessage(ctx, nil, msg)
		messageIDs = append(messageIDs, msg.MessageID)
	}
	
	successCount, failedIDs := monitor.BatchRetryMessages(ctx, messageIDs)
	assert.GreaterOrEqual(t, successCount, 0)
	assert.NotNil(t, failedIDs)
}
