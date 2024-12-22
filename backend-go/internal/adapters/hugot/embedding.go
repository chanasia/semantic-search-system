// internal/adapters/hugot/embedding.go
package hugot

import (
	"fmt"

	"github.com/knights-analytics/hugot/pipelines"
)

type HugotEmbeddingAdapter struct {
	pipeline *pipelines.FeatureExtractionPipeline
}

func NewHugotEmbeddingAdapter(pipeline interface{}) (*HugotEmbeddingAdapter, error) {
	hugotPipeline, ok := pipeline.(*pipelines.FeatureExtractionPipeline)
	if !ok {
		return nil, fmt.Errorf("invalid pipeline type")
	}
	return &HugotEmbeddingAdapter{pipeline: hugotPipeline}, nil
}

func (a *HugotEmbeddingAdapter) GenerateEmbedding(text string) (*[]float32, error) {
	result, err := a.pipeline.RunPipeline([]string{text})
	if err != nil {
		return nil, err
	}

	embedding := result.Embeddings[0]
	return &embedding, nil
}

func (a *HugotEmbeddingAdapter) GenerateEmbeddings(texts []string) (*[][]float32, error) {
	result, err := a.pipeline.RunPipeline(texts)
	if err != nil {
		return nil, err
	}

	return &result.Embeddings, nil
}
