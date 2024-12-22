package services

import (
	"github.com/chanasia/semantic-search-system/internal/core/ports"
)

type EmbeddingService struct {
	embeddingAdapter ports.EmbeddingService
	modelDownloader  ports.ModelDownloader
}

func NewEmbeddingService(
	embeddingAdapter ports.EmbeddingService,
	modelDownloader ports.ModelDownloader,
) *EmbeddingService {
	return &EmbeddingService{
		embeddingAdapter: embeddingAdapter,
		modelDownloader:  modelDownloader,
	}
}

func (s *EmbeddingService) GenerateEmbedding(text string) (*[]float32, error) {
	return s.embeddingAdapter.GenerateEmbedding(text)
}

func (s *EmbeddingService) GenerateEmbeddings(texts []string) (*[][]float32, error) {
	return s.embeddingAdapter.GenerateEmbeddings(texts)
}
