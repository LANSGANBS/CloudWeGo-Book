package mq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type StockConsumer struct{}

func NewStockConsumer() *StockConsumer {
	return &StockConsumer{}
}

func (c *StockConsumer) HandleMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		stockMsg, err := ParseStockMessage(msg.Body)
		if err != nil {
			klog.Errorf("Failed to parse message: %v", err)
			return consumer.ConsumeRetryLater, nil
		}

		klog.Infof("Processing stock message: msgId=%s, productId=%d, quantity=%d, operation=%s",
			stockMsg.MessageId, stockMsg.ProductId, stockMsg.Quantity, stockMsg.Operation)

		var processErr error
		switch stockMsg.Operation {
		case OperationDeduct:
			processErr = c.processDeductWithOrder(ctx, stockMsg)
		case OperationRestore:
			processErr = c.processRestoreWithOrder(ctx, stockMsg)
		default:
			klog.Errorf("Unknown operation: %s", stockMsg.Operation)
			continue
		}

		if processErr != nil {
			klog.Errorf("Failed to process message: msgId=%s, err=%v", stockMsg.MessageId, processErr)

			if errors.Is(processErr, ErrDuplicateMessage) {
				klog.Infof("Duplicate message, skip: msgId=%s", stockMsg.MessageId)
				continue
			}

			if stockMsg.RetryCount >= 3 {
				klog.Errorf("Message retry count exceeded, send to DLQ: msgId=%s", stockMsg.MessageId)
				c.saveToDLQ(stockMsg, processErr.Error())
				continue
			}

			return consumer.ConsumeRetryLater, nil
		}

		klog.Infof("Successfully processed message: msgId=%s", stockMsg.MessageId)
	}

	return consumer.ConsumeSuccess, nil
}

var ErrDuplicateMessage = errors.New("duplicate message")

func (c *StockConsumer) processDeductWithOrder(ctx context.Context, msg *StockMessage) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		var processed model.StockMessageLog
		err := tx.Where("message_id = ?", msg.MessageId).First(&processed).Error
		if err == nil {
			return ErrDuplicateMessage
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var stock model.Stock
		err = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("product_id = ?", msg.ProductId).
			First(&stock).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("stock record not found for product %d", msg.ProductId)
			}
			return err
		}

		if stock.Available < msg.Quantity {
			return fmt.Errorf("insufficient stock: available=%d, required=%d", stock.Available, msg.Quantity)
		}

		beforeQty := stock.Quantity
		afterQty := beforeQty - msg.Quantity

		err = tx.Model(&stock).Updates(map[string]interface{}{
			"quantity":  afterQty,
			"available": gorm.Expr("available - ?", msg.Quantity),
		}).Error
		if err != nil {
			return err
		}

		stockLog := &model.StockLog{
			ProductId:    msg.ProductId,
			OrderNo:      msg.OrderNo,
			ChangeType:   model.ChangeTypeSale,
			ChangeQty:    -msg.Quantity,
			BeforeQty:    beforeQty,
			AfterQty:     afterQty,
			OperatorId:   msg.UserId,
			OperatorName: "system",
			Remark:       fmt.Sprintf("RocketMQ deduct, msgId=%s", msg.MessageId),
		}

		if err := tx.Create(stockLog).Error; err != nil {
			return err
		}

		msgLog := &model.StockMessageLog{
			MessageId:   msg.MessageId,
			ProductId:   msg.ProductId,
			Operation:   msg.Operation,
			Quantity:    msg.Quantity,
			Status:      "processed",
			ProcessedAt: time.Now(),
		}
		return tx.Create(msgLog).Error
	})
}

func (c *StockConsumer) processRestoreWithOrder(ctx context.Context, msg *StockMessage) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		var processed model.StockMessageLog
		err := tx.Where("message_id = ?", msg.MessageId).First(&processed).Error
		if err == nil {
			return ErrDuplicateMessage
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var stock model.Stock
		err = tx.Set("gorm:query_option", "FOR UPDATE").
			Where("product_id = ?", msg.ProductId).
			First(&stock).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stock = model.Stock{
					ProductId:   msg.ProductId,
					Quantity:    msg.Quantity,
					Available:   msg.Quantity,
					Reserved:    0,
					MinStock:    10,
					MaxStock:    1000,
					SafetyStock: 20,
					Status:      model.StockStatusNormal,
				}
				if err := tx.Create(&stock).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			beforeQty := stock.Quantity
			afterQty := beforeQty + msg.Quantity

			err = tx.Model(&stock).Updates(map[string]interface{}{
				"quantity":  afterQty,
				"available": gorm.Expr("available + ?", msg.Quantity),
			}).Error
			if err != nil {
				return err
			}

			stockLog := &model.StockLog{
				ProductId:    msg.ProductId,
				OrderNo:      msg.OrderNo,
				ChangeType:   model.ChangeTypeReturn,
				ChangeQty:    msg.Quantity,
				BeforeQty:    beforeQty,
				AfterQty:     afterQty,
				OperatorId:   msg.UserId,
				OperatorName: "system",
				Remark:       fmt.Sprintf("RocketMQ restore, msgId=%s", msg.MessageId),
			}

			if err := tx.Create(stockLog).Error; err != nil {
				return err
			}
		}

		msgLog := &model.StockMessageLog{
			MessageId:   msg.MessageId,
			ProductId:   msg.ProductId,
			Operation:   msg.Operation,
			Quantity:    msg.Quantity,
			Status:      "processed",
			ProcessedAt: time.Now(),
		}
		return tx.Create(msgLog).Error
	})
}

func (c *StockConsumer) saveToDLQ(msg *StockMessage, errMsg string) {
	dlq := &model.StockDLQ{
		MessageId:    msg.MessageId,
		ProductId:    msg.ProductId,
		Quantity:     msg.Quantity,
		OrderNo:      msg.OrderNo,
		UserId:       msg.UserId,
		Operation:    msg.Operation,
		RetryCount:   msg.RetryCount,
		ErrorMessage: errMsg,
	}

	mysql.DB.Create(dlq)
}
