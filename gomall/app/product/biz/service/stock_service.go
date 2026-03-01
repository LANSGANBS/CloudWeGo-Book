package service

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	product "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type StockService struct {
	ctx context.Context
}

func NewStockService(ctx context.Context) *StockService {
	return &StockService{ctx: ctx}
}

func (s *StockService) GetStock(productId uint32) (*model.Stock, error) {
	var stock model.Stock
	err := mysql.DB.WithContext(s.ctx).Where("product_id = ?", productId).First(&stock).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (s *StockService) InitStock(productId uint32, quantity int64, minStock, maxStock, safetyStock int64) error {
	stock := &model.Stock{
		ProductId:   productId,
		Quantity:    quantity,
		Reserved:    0,
		Available:   quantity,
		MinStock:    minStock,
		MaxStock:    maxStock,
		SafetyStock: safetyStock,
		Status:      model.StockStatusNormal,
	}
	
	err := mysql.DB.WithContext(s.ctx).Create(stock).Error
	if err != nil {
		return err
	}
	
	err = redis.InitStock(s.ctx, productId, quantity)
	if err != nil {
		klog.Errorf("Failed to init redis stock: %v", err)
	}
	
	return nil
}

func (s *StockService) UpdateStock(productId uint32, quantity int64) error {
	return mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Where("product_id = ?", productId).
		Updates(map[string]interface{}{
			"quantity":  quantity,
			"available": gorm.Expr("quantity - reserved"),
		}).Error
}

func (s *StockService) ReserveStock(productId uint32, quantity int64) error {
	return mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Where("product_id = ? AND available >= ?", productId, quantity).
		Updates(map[string]interface{}{
			"reserved":  gorm.Expr("reserved + ?", quantity),
			"available": gorm.Expr("available - ?", quantity),
		}).Error
}

func (s *StockService) ReleaseReservedStock(productId uint32, quantity int64) error {
	return mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Where("product_id = ? AND reserved >= ?", productId, quantity).
		Updates(map[string]interface{}{
			"reserved":  gorm.Expr("reserved - ?", quantity),
			"available": gorm.Expr("available + ?", quantity),
		}).Error
}

func (s *StockService) DeductStockWithLog(productId uint32, quantity int64, changeType int8, orderNo string, operatorId uint32, operatorName string) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.Stock
		err := tx.Where("product_id = ?", productId).First(&stock).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("stock record not found, please initialize stock first")
			}
			return err
		}
		
		if stock.Available < quantity {
			return errors.New("insufficient stock")
		}
		
		beforeQty := stock.Quantity
		afterQty := beforeQty - quantity
		
		err = tx.Model(&stock).Updates(map[string]interface{}{
			"quantity":  afterQty,
			"available": gorm.Expr("available - ?", quantity),
		}).Error
		if err != nil {
			return err
		}
		
		stockLog := &model.StockLog{
			ProductId:    productId,
			OrderNo:      orderNo,
			ChangeType:   changeType,
			ChangeQty:    -quantity,
			BeforeQty:    beforeQty,
			AfterQty:     afterQty,
			OperatorId:   operatorId,
			OperatorName: operatorName,
		}
		
		if err := tx.Create(stockLog).Error; err != nil {
			return err
		}
		
		if err := s.checkAndCreateAlert(tx, &stock, afterQty); err != nil {
			klog.Errorf("Failed to check alert: %v", err)
		}
		
		return nil
	})
}

func (s *StockService) AddStockWithLog(productId uint32, quantity int64, changeType int8, orderNo string, operatorId uint32, operatorName string) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.Stock
		err := tx.Where("product_id = ?", productId).First(&stock).Error
		if err != nil {
			return err
		}
		
		beforeQty := stock.Quantity
		afterQty := beforeQty + quantity
		
		err = tx.Model(&stock).Updates(map[string]interface{}{
			"quantity":  afterQty,
			"available": gorm.Expr("available + ?", quantity),
		}).Error
		if err != nil {
			return err
		}
		
		stockLog := &model.StockLog{
			ProductId:    productId,
			OrderNo:      orderNo,
			ChangeType:   changeType,
			ChangeQty:    quantity,
			BeforeQty:    beforeQty,
			AfterQty:     afterQty,
			OperatorId:   operatorId,
			OperatorName: operatorName,
		}
		
		if err := tx.Create(stockLog).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func (s *StockService) checkAndCreateAlert(tx *gorm.DB, stock *model.Stock, currentQty int64) error {
	if currentQty <= stock.SafetyStock {
		alert := &model.StockAlert{
			ProductId:    stock.ProductId,
			AlertType:    model.AlertTypeLowStock,
			AlertLevel:   model.AlertLevelDanger,
			Threshold:    stock.SafetyStock,
			CurrentValue: currentQty,
			Status:       model.AlertStatusPending,
		}
		
		if currentQty <= stock.MinStock {
			alert.AlertLevel = model.AlertLevelDanger
		} else {
			alert.AlertLevel = model.AlertLevelWarning
		}
		
		var existingAlert model.StockAlert
		err := tx.Where("product_id = ? AND alert_type = ? AND status = ?", 
			stock.ProductId, model.AlertTypeLowStock, model.AlertStatusPending).
			First(&existingAlert).Error
		if err == gorm.ErrRecordNotFound {
			return tx.Create(alert).Error
		}
	}
	
	if currentQty >= stock.MaxStock {
		alert := &model.StockAlert{
			ProductId:    stock.ProductId,
			AlertType:    model.AlertTypeOverStock,
			AlertLevel:   model.AlertLevelWarning,
			Threshold:    stock.MaxStock,
			CurrentValue: currentQty,
			Status:       model.AlertStatusPending,
		}
		
		return tx.Create(alert).Error
	}
	
	return nil
}

func (s *StockService) GetStockLogs(productId uint32, page, pageSize int32) ([]*model.StockLog, int64, error) {
	var logs []*model.StockLog
	var total int64
	
	db := mysql.DB.WithContext(s.ctx).Model(&model.StockLog{})
	if productId > 0 {
		db = db.Where("product_id = ?", productId)
	}
	
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * pageSize
	if err := db.Order("id DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	
	return logs, total, nil
}

func (s *StockService) GetPendingAlerts() ([]*model.StockAlert, error) {
	var alerts []*model.StockAlert
	err := mysql.DB.WithContext(s.ctx).
		Where("status = ?", model.AlertStatusPending).
		Order("alert_level DESC, created_at DESC").
		Find(&alerts).Error
	return alerts, err
}

func (s *StockService) HandleAlert(alertId uint32, handlerId uint32, handlerName string, status int8, remark string) error {
	now := time.Now()
	return mysql.DB.WithContext(s.ctx).Model(&model.StockAlert{}).
		Where("id = ?", alertId).
		Updates(map[string]interface{}{
			"status":       status,
			"handled_at":   &now,
			"handler_id":   handlerId,
			"handler_name": handlerName,
			"remark":       remark,
		}).Error
}

func (s *StockService) CreateStockCheck(warehouseId, operatorId uint32, operatorName string, productIds []uint32) (*model.StockCheck, error) {
	checkNo := time.Now().Format("CK20060102150405")
	
	stockCheck := &model.StockCheck{
		CheckNo:      checkNo,
		WarehouseId:  warehouseId,
		Status:       model.CheckStatusPending,
		TotalItems:   len(productIds),
		OperatorId:   operatorId,
		OperatorName: operatorName,
	}
	
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(stockCheck).Error; err != nil {
			return err
		}
		
		for _, productId := range productIds {
			var stock model.Stock
			if err := tx.Where("product_id = ?", productId).First(&stock).Error; err != nil {
				continue
			}
			
			item := &model.StockCheckItem{
				CheckId:   uint32(stockCheck.ID),
				ProductId: productId,
				SystemQty: stock.Quantity,
				ActualQty: stock.Quantity,
				DiffQty:   0,
			}
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
	
	return stockCheck, err
}

func (s *StockService) UpdateCheckItem(itemId uint32, actualQty int64, remark string) error {
	var item model.StockCheckItem
	if err := mysql.DB.WithContext(s.ctx).First(&item, itemId).Error; err != nil {
		return err
	}
	
	diffQty := actualQty - item.SystemQty
	return mysql.DB.WithContext(s.ctx).Model(&item).Updates(map[string]interface{}{
		"actual_qty": actualQty,
		"diff_qty":   diffQty,
		"remark":     remark,
	}).Error
}

func (s *StockService) FinishStockCheck(checkId uint32, operatorId uint32, operatorName string) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		var check model.StockCheck
		if err := tx.First(&check, checkId).Error; err != nil {
			return err
		}
		
		var items []model.StockCheckItem
		if err := tx.Where("check_id = ?", checkId).Find(&items).Error; err != nil {
			return err
		}
		
		diffCount := 0
		for _, item := range items {
			if item.DiffQty != 0 {
				diffCount++
				
				var stock model.Stock
				if err := tx.Where("product_id = ?", item.ProductId).First(&stock).Error; err != nil {
					continue
				}
				
				beforeQty := stock.Quantity
				afterQty := item.ActualQty
				
				tx.Model(&stock).Updates(map[string]interface{}{
					"quantity":  afterQty,
					"available": gorm.Expr("quantity - reserved"),
				})
				
				stockLog := &model.StockLog{
					ProductId:    item.ProductId,
					OrderNo:      check.CheckNo,
					ChangeType:   model.ChangeTypeCheck,
					ChangeQty:    item.DiffQty,
					BeforeQty:    beforeQty,
					AfterQty:     afterQty,
					OperatorId:   operatorId,
					OperatorName: operatorName,
					Remark:       "库存盘点调整",
				}
				tx.Create(stockLog)
			}
		}
		
		now := time.Now()
		return tx.Model(&check).Updates(map[string]interface{}{
			"status":       model.CheckStatusFinished,
			"diff_items":   diffCount,
			"finished_at":  &now,
		}).Error
	})
}

func (s *StockService) GetStockReport() (map[string]interface{}, error) {
	var totalProducts int64
	var lowStockCount int64
	var overStockCount int64
	var totalQuantity int64
	
	mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).Count(&totalProducts)
	mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Where("quantity <= safety_stock").Count(&lowStockCount)
	mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Where("quantity >= max_stock").Count(&overStockCount)
	mysql.DB.WithContext(s.ctx).Model(&model.Stock{}).
		Select("COALESCE(SUM(quantity), 0)").Scan(&totalQuantity)
	
	pendingAlerts, _ := s.GetPendingAlerts()
	
	return map[string]interface{}{
		"total_products":    totalProducts,
		"low_stock_count":   lowStockCount,
		"over_stock_count":  overStockCount,
		"total_quantity":    totalQuantity,
		"pending_alerts":    len(pendingAlerts),
		"generated_at":      time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *StockService) SyncStockToRedis(productId uint32) error {
	stock, err := s.GetStock(productId)
	if err != nil {
		return err
	}
	if stock == nil {
		return errors.New("stock not found")
	}
	return redis.InitStock(s.ctx, productId, stock.Quantity)
}

func (s *StockService) SyncAllStockToRedis() error {
	var stocks []model.Stock
	if err := mysql.DB.WithContext(s.ctx).Find(&stocks).Error; err != nil {
		return err
	}
	
	for _, stock := range stocks {
		if err := redis.InitStock(s.ctx, stock.ProductId, stock.Quantity); err != nil {
			klog.Errorf("Failed to sync stock %d to redis: %v", stock.ProductId, err)
		}
	}
	
	return nil
}

func (s *StockService) DeductStockFromRedis(req *product.DeductStockReq) (*product.DeductStockResp, error) {
	resp := &product.DeductStockResp{}
	
	remainingStock, err := redis.DeductStock(s.ctx, req.ProductId, int64(req.Quantity))
	if err != nil {
		resp.Success = false
		resp.ErrorMessage = err.Error()
		return resp, nil
	}
	
	if remainingStock < 0 {
		switch remainingStock {
		case -1:
			resp.Success = false
			resp.ErrorMessage = "insufficient stock"
		case -2:
			resp.Success = false
			resp.ErrorMessage = "stock not initialized"
		default:
			resp.Success = false
			resp.ErrorMessage = "unknown error"
		}
		return resp, nil
	}
	
	resp.Success = true
	resp.RemainingStock = remainingStock
	
	klog.Infof("Deducted stock from redis: productId=%d, quantity=%d, remaining=%d", 
		req.ProductId, req.Quantity, remainingStock)
	
	return resp, nil
}
