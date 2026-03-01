package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	VectorKeyPrefix = "product_vector"
	VectorIndexKey  = "product_vector_index"
)

type VectorStore struct{}

type ProductVector struct {
	ProductID   uint32    `json:"product_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Vector      []float32 `json:"vector"`
	Stock       int64     `json:"stock"`
	Sales       int       `json:"sales"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewVectorStore() *VectorStore {
	return &VectorStore{}
}

func (v *VectorStore) StoreProductVector(ctx context.Context, pv *ProductVector) error {
	key := fmt.Sprintf("%s:%d", VectorKeyPrefix, pv.ProductID)
	data, err := json.Marshal(pv)
	if err != nil {
		return fmt.Errorf("failed to marshal product vector: %w", err)
	}

	if err := redis.RedisClient.Set(ctx, key, data, 0).Err(); err != nil {
		return fmt.Errorf("failed to store product vector: %w", err)
	}

	klog.Infof("Stored vector for product %d", pv.ProductID)
	return nil
}

func (v *VectorStore) GetAllProductVectors(ctx context.Context) ([]*ProductVector, error) {
	pattern := fmt.Sprintf("%s:*", VectorKeyPrefix)
	keys, err := redis.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get vector keys: %w", err)
	}

	var vectors []*ProductVector
	for _, key := range keys {
		data, err := redis.RedisClient.Get(ctx, key).Bytes()
		if err != nil {
			klog.Warnf("Failed to get vector for key %s: %v", key, err)
			continue
		}

		var pv ProductVector
		if err := json.Unmarshal(data, &pv); err != nil {
			klog.Warnf("Failed to unmarshal vector for key %s: %v", key, err)
			continue
		}
		vectors = append(vectors, &pv)
	}

	klog.Infof("Retrieved %d product vectors", len(vectors))
	return vectors, nil
}

func (v *VectorStore) DeleteAllVectors(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", VectorKeyPrefix)
	keys, err := redis.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get vector keys: %w", err)
	}

	if len(keys) > 0 {
		if err := redis.RedisClient.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete vectors: %w", err)
		}
		klog.Infof("Deleted %d vectors", len(keys))
	}
	return nil
}

func CosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

type SearchResult struct {
	ProductID   uint32
	Name        string
	Description string
	Score       float32
}

func (v *VectorStore) calculateWeight(stock int64, sales int64, price float32, createdAt time.Time) float32 {
	var weight float32 = 1.0

	if stock <= 0 {
		weight *= 0.5
	} else if stock < 10 {
		weight *= 0.8
	}

	if sales > 1000 {
		weight += 1.0
	} else if sales > 500 {
		weight += 0.8
	} else if sales > 100 {
		weight += 0.5
	} else if sales > 50 {
		weight += 0.3
	} else if sales > 10 {
		weight += 0.15
	}

	if !createdAt.IsZero() {
		daysSinceCreated := time.Since(createdAt).Hours() / 24
		if daysSinceCreated < 7 {
			weight += 0.2
		} else if daysSinceCreated < 30 {
			weight += 0.1
		}
	}

	if price > 0 && price < 50 {
		weight += 0.1
	}

	if weight < 0.2 {
		weight = 0.2
	}

	return weight
}

func (v *VectorStore) SearchSimilar(ctx context.Context, queryVector []float32, topK int) ([]*SearchResult, error) {
	vectors, err := v.GetAllProductVectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get product vectors: %w", err)
	}

	var results []*SearchResult
	for _, pv := range vectors {
		similarity := CosineSimilarity(queryVector, pv.Vector)

		var stock int64 = 999
		var sales int64 = int64(pv.Sales)
		var price float32 = pv.Price
		var createdAt time.Time = pv.CreatedAt

		dbProduct, err := model.GetProductById(mysql.DB, ctx, int(pv.ProductID))
		if err == nil {
			sales = dbProduct.Sales
			price = dbProduct.Price
			createdAt = dbProduct.CreatedAt

			var stockData model.Stock
			if err := mysql.DB.WithContext(ctx).Where("product_id = ?", pv.ProductID).First(&stockData).Error; err == nil {
				stock = stockData.Available
			}
		} else {
			klog.Warnf("Failed to get product %d from DB, using cached values: %v", pv.ProductID, err)
			stock = pv.Stock
		}

		weight := v.calculateWeight(stock, sales, price, createdAt)
		adjustedScore := similarity * weight

		klog.Infof("Product %d: similarity=%.4f, weight=%.4f, adjustedScore=%.4f (stock=%d, sales=%d)",
			pv.ProductID, similarity, weight, adjustedScore, stock, sales)

		results = append(results, &SearchResult{
			ProductID:   pv.ProductID,
			Name:        pv.Name,
			Description: pv.Description,
			Score:       adjustedScore,
		})
	}

	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	if topK > len(results) {
		topK = len(results)
	}

	klog.Infof("Found %d similar products, returning top %d", len(results), topK)
	return results[:topK], nil
}
