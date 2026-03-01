package message

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/stretchr/testify/assert"
)

func TestParsePayload(t *testing.T) {
	payload := `{"product_id": 123, "quantity": 10, "order_no": "order-001"}`
	
	result, err := ParsePayload[StockDeductPayload](payload)
	assert.NoError(t, err)
	assert.Equal(t, uint32(123), result.ProductID)
	assert.Equal(t, int32(10), result.Quantity)
	assert.Equal(t, "order-001", result.OrderNo)
}

func TestGenerateMessageID(t *testing.T) {
	id1 := generateMessageID(model.MessageTypeStockDeduct)
	id2 := generateMessageID(model.MessageTypeStockDeduct)
	
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Contains(t, id1, "msg_stock_deduct_")
	assert.Contains(t, id2, "msg_stock_deduct_")
}

func TestCalculateNextRetry(t *testing.T) {
	now := time.Now()
	next := calculateNextRetry(now)
	
	assert.True(t, next.After(now))
	assert.Equal(t, 2*time.Minute, next.Sub(now).Round(time.Second))
}

func TestLocalMessageService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	service := GetMessageService()
	if service == nil {
		t.Skip("Message service not initialized")
	}
	
	ctx := context.Background()
	
	t.Run("CreateMessage", func(t *testing.T) {
		msg := &model.LocalMessage{
			MessageID:     "test-msg-integration-001",
			BusinessID:    "order-integration-001",
			MessageType:   model.MessageTypeStockDeduct,
			TargetService: "product",
			TargetMethod:  "DeductStock",
			Payload:       `{"product_id": 1, "quantity": 10}`,
		}
		
		err := service.CreateMessage(ctx, nil, msg)
		assert.NoError(t, err)
		
		saved, err := service.GetMessageByID(ctx, msg.MessageID)
		assert.NoError(t, err)
		assert.Equal(t, model.MessageStatusPending, saved.Status)
		assert.Equal(t, 5, saved.MaxRetry)
	})
	
	t.Run("UpdateStatus", func(t *testing.T) {
		msgID := "test-msg-integration-002"
		msg := &model.LocalMessage{
			MessageID:     msgID,
			BusinessID:    "order-integration-002",
			MessageType:   model.MessageTypeStockDeduct,
			TargetService: "product",
			TargetMethod:  "DeductStock",
			Payload:       `{}`,
			Status:        model.MessageStatusPending,
		}
		service.CreateMessage(ctx, nil, msg)
		
		err := service.UpdateStatus(ctx, msgID, model.MessageStatusConfirmed, "")
		assert.NoError(t, err)
		
		updated, err := service.GetMessageByID(ctx, msgID)
		assert.NoError(t, err)
		assert.Equal(t, model.MessageStatusConfirmed, updated.Status)
		assert.NotNil(t, updated.ConfirmAt)
	})
	
	t.Run("GetPendingMessages", func(t *testing.T) {
		messages, err := service.GetPendingMessages(ctx, 10)
		assert.NoError(t, err)
		assert.NotNil(t, messages)
	})
	
	t.Run("GetMessageStats", func(t *testing.T) {
		stats, err := service.GetMessageStats(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, stats)
	})
}
