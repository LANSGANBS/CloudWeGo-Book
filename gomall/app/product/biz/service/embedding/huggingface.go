package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SiliconFlowClient struct {
	apiToken   string
	model      string
	baseURL    string
	httpClient *http.Client
}

func NewSiliconFlowClient(apiToken string) *SiliconFlowClient {
	return &SiliconFlowClient{
		apiToken: apiToken,
		model:    "BAAI/bge-large-zh-v1.5",
		baseURL:  "https://api.siliconflow.cn/v1/embeddings",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
		Object    string    `json:"object"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *SiliconFlowClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody := map[string]interface{}{
		"model":           c.model,
		"input":           text,
		"encoding_format": "float",
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		klog.Errorf("SiliconFlow API error: status=%d, body=%s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("SiliconFlow API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result EmbeddingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	klog.Infof("Got embedding for text '%s': dimension=%d", text, len(result.Data[0].Embedding))
	return result.Data[0].Embedding, nil
}

func (c *SiliconFlowClient) GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	results := make([][]float32, len(texts))
	for i, text := range texts {
		embedding, err := c.GetEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to get embedding for text %d: %w", i, err)
		}
		results[i] = embedding
	}
	return results, nil
}

type HuggingFaceClient = SiliconFlowClient

func NewHuggingFaceClient(apiToken, _ string) *SiliconFlowClient {
	return NewSiliconFlowClient(apiToken)
}
