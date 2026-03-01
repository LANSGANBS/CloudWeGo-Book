package redis

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

//go:embed lua/deduct_stock.lua
var deductStockScript string

var deductStockSHA string

func InitLuaScripts() error {
	var err error
	deductStockSHA, err = RedisClient.ScriptLoad(context.Background(), deductStockScript).Result()
	if err != nil {
		return fmt.Errorf("failed to load deduct_stock script: %w", err)
	}
	klog.Infof("Loaded Lua script deduct_stock with SHA: %s", deductStockSHA)
	return nil
}

func GetStockKey(productId uint32) string {
	return fmt.Sprintf("stock:product:%d", productId)
}

func InitStock(ctx context.Context, productId uint32, stock int64) error {
	key := GetStockKey(productId)
	return RedisClient.Set(ctx, key, stock, 24*time.Hour).Err()
}

func GetStock(ctx context.Context, productId uint32) (int64, error) {
	key := GetStockKey(productId)
	result, err := RedisClient.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func DeductStock(ctx context.Context, productId uint32, quantity int64) (int64, error) {
	key := GetStockKey(productId)
	
	result, err := RedisClient.EvalSha(ctx, deductStockSHA, []string{key}, quantity, time.Now().Unix()).Int64()
	if err != nil {
		if err.Error() == "NOSCRIPT No matching script. Please use EVAL." {
			deductStockSHA, err = RedisClient.ScriptLoad(ctx, deductStockScript).Result()
			if err != nil {
				return 0, fmt.Errorf("failed to reload script: %w", err)
			}
			return DeductStock(ctx, productId, quantity)
		}
		return 0, fmt.Errorf("failed to execute deduct_stock script: %w", err)
	}
	
	return result, nil
}

func RestoreStock(ctx context.Context, productId uint32, quantity int64) error {
	key := GetStockKey(productId)
	return RedisClient.IncrBy(ctx, key, quantity).Err()
}

func RefreshStockTTL(ctx context.Context, productId uint32) error {
	key := GetStockKey(productId)
	return RedisClient.Expire(ctx, key, 24*time.Hour).Err()
}
