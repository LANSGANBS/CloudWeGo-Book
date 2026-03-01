package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/kitex/pkg/klog"
)

type StockMessage struct {
	MessageId    string    `json:"message_id"`
	ProductId    uint32    `json:"product_id"`
	Quantity     int64     `json:"quantity"`
	OrderNo      string    `json:"order_no"`
	UserId       uint32    `json:"user_id"`
	Operation    string    `json:"operation"`
	Timestamp    int64     `json:"timestamp"`
	RetryCount   int       `json:"retry_count"`
}

const (
	OperationDeduct  = "deduct"
	OperationRestore = "restore"
)

func NewStockMessage(productId uint32, quantity int64, orderNo string, userId uint32, operation string) *StockMessage {
	return &StockMessage{
		MessageId:  fmt.Sprintf("stock_%d_%d", productId, time.Now().UnixNano()),
		ProductId:  productId,
		Quantity:   quantity,
		OrderNo:    orderNo,
		UserId:     userId,
		Operation:  operation,
		Timestamp:  time.Now().Unix(),
		RetryCount: 0,
	}
}

func SendStockDeductMessage(ctx context.Context, msg *StockMessage) error {
	if ProducerInstance == nil {
		return fmt.Errorf("rocketmq producer not initialized")
	}
	
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	m := &primitive.Message{
		Topic: TopicStockDeduct,
		Body:  body,
	}
	
	m.WithShardingKey(fmt.Sprintf("%d", msg.ProductId))
	m.WithTag(msg.Operation)
	m.WithKeys([]string{msg.MessageId})
	
	result, err := ProducerInstance.SendSync(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	
	klog.Infof("Sent stock deduct message: productId=%d, quantity=%d, msgId=%s, result=%s", 
		msg.ProductId, msg.Quantity, msg.MessageId, result.String())
	
	return nil
}

func SendStockRestoreMessage(ctx context.Context, msg *StockMessage) error {
	if ProducerInstance == nil {
		return fmt.Errorf("rocketmq producer not initialized")
	}
	
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	
	m := &primitive.Message{
		Topic: TopicStockRestore,
		Body:  body,
	}
	
	m.WithShardingKey(fmt.Sprintf("%d", msg.ProductId))
	m.WithTag(msg.Operation)
	m.WithKeys([]string{msg.MessageId})
	
	result, err := ProducerInstance.SendSync(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	
	klog.Infof("Sent stock restore message: productId=%d, quantity=%d, msgId=%s, result=%s", 
		msg.ProductId, msg.Quantity, msg.MessageId, result.String())
	
	return nil
}

func ParseStockMessage(body []byte) (*StockMessage, error) {
	var msg StockMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return &msg, nil
}
