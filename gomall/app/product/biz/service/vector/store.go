package vector

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/redis"
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

func (v *VectorStore) calculateWeight(pv *ProductVector) float32 {
	var weight float32 = 1.0

	if pv.Stock <= 0 {
		weight *= 0.3
	} else if pv.Stock < 10 {
		weight *= 0.7
	}

	if pv.Sales > 1000 {
		weight += 0.3
	} else if pv.Sales > 100 {
		weight += 0.2
	} else if pv.Sales > 10 {
		weight += 0.1
	}

	daysSinceCreated := time.Since(pv.CreatedAt).Hours() / 24
	if daysSinceCreated < 7 {
		weight += 0.2
	} else if daysSinceCreated < 30 {
		weight += 0.1
	}

	if pv.Price > 0 && pv.Price < 50 {
		weight += 0.1
	}

	if weight < 0.1 {
		weight = 0.1
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
		weight := v.calculateWeight(pv)
		adjustedScore := similarity * weight
		
		klog.Infof("Product %d: similarity=%.4f, weight=%.4f, adjustedScore=%.4f (stock=%d, sales=%d)", 
			pv.ProductID, similarity, weight, adjustedScore, pv.Stock, pv.Sales)
		
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
