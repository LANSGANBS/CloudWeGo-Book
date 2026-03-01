package message

import (
	"context"
	"sync"

	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

var (
	messageService *LocalMessageService
	messageSender  *MessageSender
	once           sync.Once
)

func Init() {
	once.Do(func() {
		messageService = NewLocalMessageService()
		
		config := DefaultSenderConfig()
		messageSender = NewMessageSender(messageService, config)
		
		messageSender.RegisterHandler(model.MessageTypeStockDeduct, NewStockDeductHandler().Handle)
		messageSender.RegisterHandler(model.MessageTypeStockRestore, NewStockRestoreHandler().Handle)
		
		klog.Info("Message module initialized")
	})
}

func Start(ctx context.Context) {
	if messageSender != nil {
		messageSender.Start(ctx)
	}
}

func Stop() {
	if messageSender != nil {
		messageSender.Stop()
	}
}

func GetMessageService() *LocalMessageService {
	return messageService
}

func GetMessageSender() *MessageSender {
	return messageSender
}

func CreateMessageInTransaction(ctx context.Context, tx *gorm.DB, businessID string, msgType model.MessageType, targetService, targetMethod string, payload interface{}) (*model.LocalMessage, error) {
	return messageService.CreateMessageWithBusiness(ctx, tx, businessID, msgType, targetService, targetMethod, payload)
}

func CreateStockDeductMessage(ctx context.Context, tx *gorm.DB, businessID string, productID uint32, quantity int32, orderNo string, userID uint32) (*model.LocalMessage, error) {
	payload := StockDeductPayload{
		ProductID: productID,
		Quantity:  quantity,
		OrderNo:   orderNo,
		UserId:    userID,
	}
	return CreateMessageInTransaction(ctx, tx, businessID, model.MessageTypeStockDeduct, "product", "DeductStock", payload)
}

func CreateStockRestoreMessage(ctx context.Context, tx *gorm.DB, businessID string, productID uint32, quantity int32, orderNo string, userID uint32, reason string) (*model.LocalMessage, error) {
	payload := StockRestorePayload{
		ProductID: productID,
		Quantity:  quantity,
		OrderNo:   orderNo,
		UserId:    userID,
		Reason:    reason,
	}
	return CreateMessageInTransaction(ctx, tx, businessID, model.MessageTypeStockRestore, "product", "RestoreStock", payload)
}

func GetDB() *gorm.DB {
	return mysql.DB
}
