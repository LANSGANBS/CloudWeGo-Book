package rag

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/model"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service/embedding"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service/vector"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	ragService     *RAGService
	ragServiceOnce sync.Once
)

type RAGService struct {
	embeddingClient *embedding.HuggingFaceClient
	vectorStore     *vector.VectorStore
}

func GetRAGService() *RAGService {
	ragServiceOnce.Do(func() {
		ragService = &RAGService{
			embeddingClient: embedding.NewSiliconFlowClient(""),
			vectorStore:     vector.NewVectorStore(),
		}
	})
	return ragService
}

func InitRAGService(apiToken string) {
	ragServiceOnce.Do(func() {
		ragService = &RAGService{
			embeddingClient: embedding.NewSiliconFlowClient(apiToken),
			vectorStore:     vector.NewVectorStore(),
		}
	})
}

func (s *RAGService) IndexProducts(ctx context.Context) error {
	klog.Info("Starting product indexing for RAG...")

	var products []*model.Product
	if err := mysql.DB.WithContext(ctx).Find(&products).Error; err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	klog.Infof("Found %d products to index", len(products))

	for _, p := range products {
		text := fmt.Sprintf("%s %s", p.Name, p.Description)
		
		vec, err := s.embeddingClient.GetEmbedding(ctx, text)
		if err != nil {
			klog.Errorf("Failed to get embedding for product %d: %v", p.ID, err)
			continue
		}

		pv := &vector.ProductVector{
			ProductID:   uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Vector:      vec,
		}

		if err := s.vectorStore.StoreProductVector(ctx, pv); err != nil {
			klog.Errorf("Failed to store vector for product %d: %v", p.ID, err)
			continue
		}

		klog.Infof("Indexed product %d: %s", p.ID, p.Name)
	}

	klog.Info("Product indexing completed!")
	return nil
}

func (s *RAGService) Search(ctx context.Context, query string, topK int) ([]*vector.SearchResult, error) {
	klog.Infof("RAG Search: query='%s', topK=%d", query, topK)

	queryVector, err := s.embeddingClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get query embedding: %w", err)
	}

	results, err := s.vectorStore.SearchSimilar(ctx, queryVector, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar products: %w", err)
	}

	klog.Infof("RAG Search results: %d products found", len(results))
	for i, r := range results {
		klog.Infof("  %d. Product %d: %s (score: %.4f)", i+1, r.ProductID, r.Name, r.Score)
	}

	return results, nil
}

func (s *RAGService) ReindexAll(ctx context.Context) error {
	if err := s.vectorStore.DeleteAllVectors(ctx); err != nil {
		return fmt.Errorf("failed to clear existing vectors: %w", err)
	}
	return s.IndexProducts(ctx)
}
