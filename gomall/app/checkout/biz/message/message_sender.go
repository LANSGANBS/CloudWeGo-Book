package message

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
)

type MessageHandler func(ctx context.Context, msg *model.LocalMessage) error

type MessageSender struct {
	service  *LocalMessageService
	handlers map[model.MessageType]MessageHandler
	stopCh   chan struct{}
	wg       sync.WaitGroup
	config   *SenderConfig
}

type SenderConfig struct {
	BatchSize        int
	PollInterval     time.Duration
	RetryIntervals   []time.Duration
	EnableAsync      bool
	WorkerCount      int
}

func DefaultSenderConfig() *SenderConfig {
	return &SenderConfig{
		BatchSize:      100,
		PollInterval:   time.Second * 5,
		RetryIntervals: []time.Duration{time.Second * 10, time.Second * 30, time.Minute, time.Minute * 5, time.Minute * 15},
		EnableAsync:    true,
		WorkerCount:    3,
	}
}

func NewMessageSender(service *LocalMessageService, config *SenderConfig) *MessageSender {
	if config == nil {
		config = DefaultSenderConfig()
	}
	return &MessageSender{
		service:  service,
		handlers: make(map[model.MessageType]MessageHandler),
		stopCh:   make(chan struct{}),
		config:   config,
	}
}

func (s *MessageSender) RegisterHandler(msgType model.MessageType, handler MessageHandler) {
	s.handlers[msgType] = handler
}

func (s *MessageSender) Start(ctx context.Context) {
	klog.Info("Starting message sender...")
	
	for i := 0; i < s.config.WorkerCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx, i)
	}
	
	klog.Infof("Message sender started with %d workers", s.config.WorkerCount)
}

func (s *MessageSender) Stop() {
	klog.Info("Stopping message sender...")
	close(s.stopCh)
	s.wg.Wait()
	klog.Info("Message sender stopped")
}

func (s *MessageSender) worker(ctx context.Context, workerID int) {
	defer s.wg.Done()
	
	ticker := time.NewTicker(s.config.PollInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-s.stopCh:
			klog.Infof("Worker %d stopping...", workerID)
			return
		case <-ctx.Done():
			klog.Infof("Worker %d context cancelled", workerID)
			return
		case <-ticker.C:
			s.processMessages(ctx, workerID)
		}
	}
}

func (s *MessageSender) processMessages(ctx context.Context, workerID int) {
	messages, err := s.service.GetPendingMessages(ctx, s.config.BatchSize)
	if err != nil {
		klog.Errorf("Worker %d: Failed to get pending messages: %v", workerID, err)
		return
	}
	
	if len(messages) == 0 {
		return
	}
	
	klog.Infof("Worker %d: Processing %d messages", workerID, len(messages))
	
	for _, msg := range messages {
		select {
		case <-s.stopCh:
			return
		default:
			s.processSingleMessage(ctx, &msg)
		}
	}
}

func (s *MessageSender) processSingleMessage(ctx context.Context, msg *model.LocalMessage) {
	startTime := time.Now()
	
	handler, ok := s.handlers[msg.MessageType]
	if !ok {
		errMsg := fmt.Sprintf("No handler registered for message type: %s", msg.MessageType)
		klog.Errorf(errMsg)
		s.service.UpdateStatus(ctx, msg.MessageID, model.MessageStatusFailed, errMsg)
		return
	}
	
	err := handler(ctx, msg)
	duration := time.Since(startTime).Milliseconds()
	
	if err != nil {
		klog.Errorf("Failed to process message %s: %v", msg.MessageID, err)
		
		s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, false, err.Error(), duration)
		
		if msg.RetryCount+1 >= msg.MaxRetry {
			klog.Errorf("Message %s exceeded max retries, moving to DLQ", msg.MessageID)
			s.service.MoveToDLQ(ctx, msg.MessageID)
		} else {
			s.service.IncrementRetry(ctx, msg.MessageID, err.Error())
		}
	} else {
		klog.Infof("Successfully processed message %s", msg.MessageID)
		s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, true, "", duration)
		s.service.UpdateStatus(ctx, msg.MessageID, model.MessageStatusConfirmed, "")
	}
}

func (s *MessageSender) SendMessageAsync(ctx context.Context, msg *model.LocalMessage) error {
	if !s.config.EnableAsync {
		return s.processSingleMessageSync(ctx, msg)
	}
	
	handler, ok := s.handlers[msg.MessageType]
	if !ok {
		return fmt.Errorf("no handler registered for message type: %s", msg.MessageType)
	}
	
	go func() {
		startTime := time.Now()
		err := handler(ctx, msg)
		duration := time.Since(startTime).Milliseconds()
		
		if err != nil {
			klog.Errorf("Async message %s failed: %v", msg.MessageID, err)
			s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, false, err.Error(), duration)
			
			if msg.RetryCount+1 >= msg.MaxRetry {
				s.service.MoveToDLQ(ctx, msg.MessageID)
			} else {
				s.service.IncrementRetry(ctx, msg.MessageID, err.Error())
			}
		} else {
			s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, true, "", duration)
			s.service.UpdateStatus(ctx, msg.MessageID, model.MessageStatusConfirmed, "")
		}
	}()
	
	return nil
}

func (s *MessageSender) processSingleMessageSync(ctx context.Context, msg *model.LocalMessage) error {
	handler, ok := s.handlers[msg.MessageType]
	if !ok {
		return fmt.Errorf("no handler registered for message type: %s", msg.MessageType)
	}
	
	startTime := time.Now()
	err := handler(ctx, msg)
	duration := time.Since(startTime).Milliseconds()
	
	if err != nil {
		s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, false, err.Error(), duration)
		
		if msg.RetryCount+1 >= msg.MaxRetry {
			s.service.MoveToDLQ(ctx, msg.MessageID)
		} else {
			s.service.IncrementRetry(ctx, msg.MessageID, err.Error())
		}
		return err
	}
	
	s.service.LogRetry(ctx, msg.MessageID, msg.RetryCount+1, true, "", duration)
	s.service.UpdateStatus(ctx, msg.MessageID, model.MessageStatusConfirmed, "")
	return nil
}

func (s *MessageSender) RetryMessage(ctx context.Context, messageID string) error {
	msg, err := s.service.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}
	
	if msg.Status != model.MessageStatusPending && msg.Status != model.MessageStatusFailed {
		return fmt.Errorf("message status is not retryable: %s", msg.Status)
	}
	
	return s.processSingleMessageSync(ctx, msg)
}

type StockDeductPayload struct {
	ProductID uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	OrderNo   string `json:"order_no"`
	UserId    uint32 `json:"user_id"`
}

type StockRestorePayload struct {
	ProductID uint32 `json:"product_id"`
	Quantity  int32  `json:"quantity"`
	OrderNo   string `json:"order_no"`
	UserId    uint32 `json:"user_id"`
	Reason    string `json:"reason"`
}

type OrderCreatePayload struct {
	UserId       uint32          `json:"user_id"`
	OrderNo      string          `json:"order_no"`
	Items        []OrderItemInfo `json:"items"`
	TotalAmount  float32         `json:"total_amount"`
	UserCurrency string          `json:"user_currency"`
	Email        string          `json:"email"`
}

type OrderItemInfo struct {
	ProductID uint32  `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float32 `json:"price"`
}

func ParsePayload[T any](payload string) (*T, error) {
	var result T
	if err := json.Unmarshal([]byte(payload), &result); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return &result, nil
}
